package main

import (
	"fmt"
	"log"

	"github.com/digitalocean/godo"
)

func showRegions(client *godo.Client) {
	regions, _, err := client.Regions.List(&godo.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Available regions")
	for _, region := range regions {
		fmt.Printf("%s -> %s\n", region.Slug, region.Name)
	}
}
