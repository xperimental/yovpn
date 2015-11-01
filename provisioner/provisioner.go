package provisioner

import (
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

type Endpoint struct {
	Name      string `json:"name"`
	IP        string `json:"ip"`
	Config    string `json:"config"`
	DropletID int    `json:"droplet"`
}

type Provisioner struct {
	client *godo.Client
}

var ErrNoToken = fmt.Errorf("No token provided!")

func NewProvisioner(token string) (*Provisioner, error) {
	if len(token) == 0 {
		return nil, ErrNoToken
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	return &Provisioner{
		client: client,
	}, nil
}

func (p Provisioner) CreateEndpoint(region string) (endpoint Endpoint, err error) {
	log.Println("Creating SSH key...")
	sshKey, err := createPrivateKey()
	if err != nil {
		return
	}
	sshKeyFingerprint := publicFingerprint(sshKey)
	log.Printf("Fingerprint of key: %s\n", sshKeyFingerprint)

	log.Println("Uploading key...")
	doKey, err := uploadPublicKey(p.client, sshKey.PublicKey())
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
	droplet, err := createDroplet(p.client, doKey, region)
	if err != nil {
		return
	}
	log.Printf("Created droplet %s (%d)", droplet.Name, droplet.ID)

	log.Println("Waiting for droplet to be ready...")
	dropletIP, err := waitForNetwork(p.client, droplet.ID)
	if err != nil {
		return
	}
	log.Printf("Droplet is ready: %s", dropletIP)

	log.Println("Waiting for setup script to complete...")
	sshClient := waitForSetup(sshKey, dropletIP)
	defer sshClient.Close()
	log.Println("Setup complete.")

	log.Println("Reading secret...")
	secret, err := readSecret(sshClient)
	if err != nil {
		return
	}
	log.Printf("Successfully read secret.")

	log.Println("Writing configuration...")
	config, err := createConfigFile(dropletIP, secret)
	if err != nil {
		return
	}
	log.Println("Created configuration.")

	log.Println("Remove SSH key from DigitalOcean...")
	deletePublicKey(p.client, doKey)

	endpoint = Endpoint{
		Name:      droplet.Name,
		IP:        dropletIP,
		Config:    config,
		DropletID: droplet.ID,
	}
	return
}
