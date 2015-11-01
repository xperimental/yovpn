package provisioner

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/digitalocean/godo"
)

const (
	defaultName  = "yovpn"
	defaultImage = "ubuntu-14-04-x64"
	defaultSize  = "512mb"
)

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

func createDroplet(client *godo.Client, key *godo.Key, region string) (*godo.Droplet, error) {
	userData := readCloudConfig()
	createRequest := &godo.DropletCreateRequest{
		Name:   defaultName,
		Region: region,
		Size:   defaultSize,
		Image:  godo.DropletCreateImage{Slug: defaultImage},
		SSHKeys: []godo.DropletCreateSSHKey{
			godo.DropletCreateSSHKey{Fingerprint: key.Fingerprint},
		},
		Backups:           false,
		IPv6:              false,
		PrivateNetworking: false,
		UserData:          userData,
	}
	drop, _, err := client.Droplets.Create(createRequest)
	return drop, err
}

func waitForNetwork(client *godo.Client, dropletID int) (string, error) {
	for {
		drop, _, err := client.Droplets.Get(dropletID)
		if err != nil {
			return "", err
		}

		if drop.Status == "active" {
			if len(drop.Networks.V4) > 0 {
				return drop.Networks.V4[0].IPAddress, nil
			}
		}
		<-time.After(time.Second * 5)
	}
}

func deleteDroplet(client *godo.Client, id int) error {
	drop, _, err := client.Droplets.Get(id)
	if err != nil {
		return err
	}

	log.Printf("Deleting %s (%d)", drop.Name, drop.ID)
	_, err = client.Droplets.Delete(drop.ID)
	return err
}
