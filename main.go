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
	log.Printf("Created droplet %s", drop.Name)

	log.Println("Waiting for droplet to be ready...")
	for {
		drop, _, err := client.Droplets.Get(drop.ID)
		if err != nil {
			log.Fatal(err)
		}
		if drop.Status == "active" {
			break
		}
		<-time.After(time.Second)
	}
	log.Println("Droplet is ready.")
	log.Printf("Networks: %+v", drop.Networks)
}
