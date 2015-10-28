package main

import (
	"flag"
	"fmt"

	"golang.org/x/crypto/ssh"

	"github.com/digitalocean/godo"
)

var keyName = flag.String("key", "yovpn", "SSH key to use for droplet")

var errKeyNotFound = fmt.Errorf("No key found with name %s", *keyName)

func loadPublicKey(client *godo.Client) (*godo.Key, error) {
	keys, _, err := client.Keys.List(&godo.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		if key.Name == *keyName {
			return &key, nil
		}
	}
	return nil, errKeyNotFound
}

func uploadPublicKey(client *godo.Client, publicKey ssh.PublicKey) (*godo.Key, error) {
	createRequest := &godo.KeyCreateRequest{
		Name:      *keyName,
		PublicKey: string(ssh.MarshalAuthorizedKey(publicKey)),
	}
	key, _, err := client.Keys.Create(createRequest)
	return key, err
}
