package main

import (
	"flag"
	"log"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

var token = flag.String("token", "", "DigitalOcean access token")

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

	log.Println("Reading secret...")
	secret := readSecret(dropletIP)
	log.Printf("Successfully read secret.")

	log.Println("Writing configuration...")
	file := writeConfigFile(dropletIP, secret)
	log.Printf("Written configuration to %s", file)
}
