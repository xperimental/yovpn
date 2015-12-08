// Package provisioner provides an interface into DigitalOcean which can be used to provision a VPN endpoint.
package provisioner

import (
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// Provisioner can be used to provision VPN endpoints.
type Provisioner interface {
	// CreateEndpoint stats the creation of a new endpoint in the specified region.
	// The provisioning itself will happen in a goroutine which will emit a signal when finished.
	CreateEndpoint(string) Endpoint

	// GetEndpoint returns the data for an endpoint with the specified ID or an error if not found.
	GetEndpoint(string) (Endpoint, error)

	// ListEndpoints returns a slice of all known endpoints. Not all endpoints will have connection information
	// if they were created by another provisioner.
	ListEndpoints() []Endpoint

	// DestroyEndpoint removes the endpoint from DigitalOcean.
	// If the endpoint can not be found this function will return an error.
	DestroyEndpoint(string) (Endpoint, error)

	// ListRegions returns a slice with all known regions.
	ListRegions() ([]Region, error)

	// Signal returns the signal channel used by this provisioner.
	Signal() chan struct{}
}

type provisioner struct {
	client    *godo.Client
	endpoints map[string]*Endpoint
	signal    chan struct{}
}

// ErrNoToken is used when the token is empty.
var ErrNoToken = fmt.Errorf("No token provided.")

// ErrNotFound is used when the endpoint can not be found.
var ErrNotFound = fmt.Errorf("Endpoint not found!")

// ErrTokenInvalid is used when the provided token is not valid.
var ErrTokenInvalid = fmt.Errorf("Token is not valid!")

func checkToken(client *godo.Client) bool {
	log.Println("Checking token...")
	account, _, err := client.Account.Get()
	if err != nil {
		return false
	}

	return account.Status == "active"
}

// NewProvisioner creates a new Provisioner with the specified DigitalOcean token.
func NewProvisioner(token string) (Provisioner, error) {
	if len(token) == 0 {
		return nil, ErrNoToken
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)
	if !checkToken(client) {
		return nil, ErrTokenInvalid
	}

	result := &provisioner{
		client:    client,
		endpoints: make(map[string]*Endpoint),
		signal:    make(chan struct{}),
	}
	go result.restoreEndpoints()

	return result, nil
}

func (p provisioner) Signal() chan struct{} {
	return p.signal
}

func (p provisioner) provisionEndpoint(endpoint *Endpoint, region string) {
	log.Println("Creating SSH key...")
	sshKey, err := createPrivateKey()
	if err != nil {
		return
	}
	sshKeyFingerprint := publicFingerprint(sshKey)
	log.Printf("Fingerprint of key: %s\n", sshKeyFingerprint)

	log.Println("Uploading key...")
	doKey, err := uploadPublicKey(p.client, sshKey.PublicKey(), endpoint.ID)
	if err != nil {
		return
	}
	log.Println("Uploaded key.")

	if doKey.Fingerprint != sshKeyFingerprint {
		err = fmt.Errorf("Key fingerprints do not match: %s, %s", sshKeyFingerprint, doKey.Fingerprint)
		return
	}
	log.Printf("Using key with fingerprint %s", doKey.Fingerprint)

	log.Println("Creating droplet...")
	droplet, err := createDroplet(p.client, doKey, region, endpoint.ID)
	if err != nil {
		return
	}
	endpoint.DropletID = droplet.ID
	log.Printf("Created droplet %s (%d)", droplet.Name, droplet.ID)

	log.Println("Waiting for droplet to be ready...")
	endpoint.IP, err = waitForNetwork(p.client, droplet.ID)
	if err != nil {
		return
	}
	log.Printf("Droplet is ready: %s", endpoint.IP)

	log.Println("Waiting for setup script to complete...")
	sshClient := waitForSetup(sshKey, endpoint.IP)
	defer sshClient.Close()
	log.Println("Setup complete.")

	log.Println("Reading secret...")
	secret, err := readSecret(sshClient)
	if err != nil {
		return
	}
	log.Printf("Successfully read secret.")

	log.Println("Writing configuration...")
	endpoint.Config, err = createConfigFile(endpoint.IP, secret)
	if err != nil {
		return
	}
	log.Println("Created configuration.")

	log.Println("Remove SSH key from DigitalOcean...")
	deletePublicKey(p.client, doKey)

	endpoint.Status = Running
	p.signal <- struct{}{}
}

func (p provisioner) unprovisionEndpoint(endpoint *Endpoint) {
	err := deleteDroplet(p.client, endpoint.DropletID)
	if err == nil {
		endpoint.Status = Destroyed
	}
}
