package main

import (
	"io/ioutil"
	"log"

	"github.com/xperimental/yovpn/provisioner"
)

const (
	fileMode = 0600
)

func provisionEndpoint(p *provisioner.Provisioner) {
	log.Println("Creating endpoint...")
	e := p.CreateEndpoint(region)

	log.Println("Waiting for endpoint to be ready...")
	<-p.Signal

	e, err := p.GetEndpoint(e.ID)
	if err != nil {
		log.Printf("Error getting endpoint: %s", err)
	}

	if e.Status == provisioner.Running {
		log.Println("Endpoint running. Writing configuration file...")
		err := ioutil.WriteFile(configFile, []byte(e.Config), fileMode)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func destroyEndpoints(p *provisioner.Provisioner) {
	log.Println("Destroying endpoints...")
	endpoints := p.ListEndpoints()
	for _, e := range endpoints {
		p.DestroyEndpoint(e.ID)
	}
}
