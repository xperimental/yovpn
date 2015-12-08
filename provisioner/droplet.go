package provisioner

import (
	"log"
	"time"

	"github.com/digitalocean/godo"
	"github.com/xperimental/yovpn/provisioner/config"
)

const (
	baseName     = "yovpn-"
	defaultImage = "ubuntu-14-04-x64"
	defaultSize  = "512mb"
)

func createDroplet(client *godo.Client, key *godo.Key, region string, id string) (*godo.Droplet, error) {
	createRequest := &godo.DropletCreateRequest{
		Name:   baseName + id,
		Region: region,
		Size:   defaultSize,
		Image:  godo.DropletCreateImage{Slug: defaultImage},
		SSHKeys: []godo.DropletCreateSSHKey{
			godo.DropletCreateSSHKey{Fingerprint: key.Fingerprint},
		},
		Backups:           false,
		IPv6:              false,
		PrivateNetworking: false,
		UserData:          config.CloudConfig,
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
