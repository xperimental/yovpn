package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/digitalocean/godo"
)

var image = flag.String("image", "ubuntu-14-04-x64", "Default image for droplet")
var size = flag.String("size", "512mb", "Default size for droplet")
var region = flag.String("region", "nyc2", "Default region for droplet")

func readCloudConfig() string {
	file, err := os.Open("share/cloudconfig.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(buf)
}

func createDroplet(client *godo.Client, key *godo.Key) *godo.Droplet {
	userData := readCloudConfig()
	createRequest := &godo.DropletCreateRequest{
		Name:   "yovpn",
		Region: *region,
		Size:   *size,
		Image:  godo.DropletCreateImage{Slug: *image},
		SSHKeys: []godo.DropletCreateSSHKey{
			godo.DropletCreateSSHKey{Fingerprint: key.Fingerprint},
		},
		Backups:           false,
		IPv6:              false,
		PrivateNetworking: false,
		UserData:          userData,
	}
	drop, _, err := client.Droplets.Create(createRequest)
	if err != nil {
		log.Fatal(err)
	}
	return drop
}
