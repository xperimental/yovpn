package main

import (
	"flag"
	"log"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

var token = flag.String("token", "", "DigitalOcean access token")
var regions = flag.Bool("regions", false, "If true, will display all available regions")
var destroy = flag.Bool("destroy", false, "If true, will remove all droplets with given name")

func initClient() *godo.Client {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *token})
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	return godo.NewClient(oauthClient)
}

func waitForNetwork(client *godo.Client, dropletID int) string {
	for {
		drop, _, err := client.Droplets.Get(dropletID)
		if err != nil {
			log.Fatal(err)
		}
		if drop.Status == "active" {
			if len(drop.Networks.V4) > 0 {
				return drop.Networks.V4[0].IPAddress
			}
		}
		<-time.After(time.Second * 5)
	}
}

func main() {
	flag.Parse()

	if len(*token) == 0 {
		log.Fatal("Token can not be empty!")
		return
	}

	client := initClient()

	if *regions {
		log.Println("Showing regions...")
		showRegions(client)
		return
	}

	if *destroy {
		log.Println("Removing droplets...")
		deleteDroplets(client)
		return
	}

	log.Println("Loading SSH key...")
	sshKey, err := loadPrivateKey()
	switch {
	case err == errKeyNotFound:
		log.Println("SSH Key not found! Creating...")
		sshKey = createPrivateKey()
	case err != nil:
		log.Fatal(err)
	}
	sshKeyFingerprint := publicFingerprint(sshKey)
	log.Printf("Loaded key: %s\n", sshKeyFingerprint)

	log.Println("Looking for key...")
	key, err := loadPublicKey(client)
	switch {
	case err == errKeyNotFound:
		log.Println("Key not found. Uploading...")
		key, err = uploadPublicKey(client, sshKey.PublicKey())
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Uploaded key.")
	case err != nil:
		log.Fatal(err)
	}
	if key.Fingerprint != sshKeyFingerprint {
		log.Printf("Local:  %s", sshKeyFingerprint)
		log.Printf("Remote: %s", key.Fingerprint)
		log.Fatalf("Key fingerprints do not match!", sshKeyFingerprint, key.Fingerprint)
	}
	log.Printf("Using key with fingerprint %s", key.Fingerprint)

	log.Println("Creating droplet...")
	drop := createDroplet(client, key)
	log.Printf("Created droplet %s (%d)", drop.Name, drop.ID)

	log.Println("Waiting for droplet to be ready...")
	dropletIP := waitForNetwork(client, drop.ID)
	log.Printf("Droplet is ready: %s", dropletIP)

	log.Println("Waiting for setup script to complete...")
	sshClient := waitForSetup(sshKey, dropletIP)
	defer sshClient.Close()
	log.Println("Setup complete.")

	log.Println("Reading secret...")
	secret := readSecret(sshClient)
	log.Printf("Successfully read secret.")

	log.Println("Writing configuration...")
	file := writeConfigFile(dropletIP, secret)
	log.Printf("Written configuration to %s", file)
}
