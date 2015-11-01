package provisioner

import (
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

type Provisioner struct {
	client    *godo.Client
	endpoints map[string]*Endpoint
}

var (
	ErrNoToken  = fmt.Errorf("No token provided!")
	ErrNotFound = fmt.Errorf("Endpoint not found!")
)

func checkToken(client *godo.Client) bool {
	log.Println("Checking token...")
	account, _, err := client.Account.Get()
	if err != nil {
		return false
	}

	return account.Status == "active"
}

func NewProvisioner(token string) (*Provisioner, error) {
	if len(token) == 0 {
		return nil, ErrNoToken
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)
	if !checkToken(client) {
		return nil, fmt.Errorf("Token is not valid!")
	}

	return &Provisioner{
		client:    client,
		endpoints: make(map[string]*Endpoint),
	}, nil
}

func (p Provisioner) provisionEndpoint(endpoint *Endpoint, region string) {
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
}
