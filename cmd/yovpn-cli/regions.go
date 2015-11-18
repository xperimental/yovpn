package main

import (
	"log"

	"github.com/xperimental/yovpn/provisioner"
)

func listRegions(p *provisioner.Provisioner) {
	regions, err := p.ListRegions()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Available regions:")
	for _, r := range regions {
		log.Printf("%s -> %s", r.Name, r.Description)
	}
}
