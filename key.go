package main

import (
	"flag"
	"fmt"

	"github.com/digitalocean/godo"
)

var keyName = flag.String("key", "yovpn", "SSH key to use for droplet")

func createKey(client *godo.Client) (*godo.Key, error) {
	keys, _, err := client.Keys.List(&godo.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		if key.Name == *keyName {
			return &key, nil
		}
	}
	return nil, fmt.Errorf("No key found with name %s", *keyName)
}
