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
var destroy = flag.Bool("destroy", false, "If true, will remove all yovpn droplets")

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

	client := initClient()

	if *regions {
		showRegions(client)
		return
	}

	if *destroy {
		deleteDroplets(client)
		return
	}

	log.Println("Looking for key...")
	key, err := createKey(client)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Using key with fingerprint %s", key.Fingerprint)

	log.Println("Creating droplet...")
	drop := createDroplet(client, key)
	log.Printf("Created droplet %s (%d)", drop.Name, drop.ID)

	log.Println("Waiting for droplet to be ready...")
	dropletIP := waitForNetwork(client, drop.ID)
	log.Printf("Droplet is ready: %s", dropletIP)

	log.Println("Waiting for setup script to complete...")
	sshClient := waitForSetup(dropletIP)
	defer sshClient.Close()
	log.Println("Setup complete.")

	log.Println("Reading secret...")
	secret := readSecret(sshClient)
	log.Printf("Successfully read secret.")

	log.Println("Writing configuration...")
	file := writeConfigFile(dropletIP, secret)
	log.Printf("Written configuration to %s", file)
}
