package main

import (
	"log"

	"github.com/xperimental/yovpn/provisioner"
)

func main() {
	checkFlags()

	log.Println("Creating provisioner...")
	p, err := provisioner.NewProvisioner(token)
	if err != nil {
		log.Fatal(err)
	}
	<-p.Signal

	if destroy {
		destroyEndpoints(p)
	} else {
		if len(region) == 0 {
			log.Println("You need to provide a region to provision in.")
			listRegions(p)
			return
		}

		provisionEndpoint(p)
	}
}
