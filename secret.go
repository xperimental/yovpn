package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

var keyfile = flag.String("keyfile", "yovpn.key", "SSH keyfile to use for connection")

func readKey() []byte {
	file, err := os.Open(*keyfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return buf
}

func createSSHClient(ip string) (*ssh.Client, error) {
	key := readKey()
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal(err)
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	return ssh.Dial("tcp", net.JoinHostPort(ip, "22"), config)
}

func readSecret(client *ssh.Client) string {
	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("cat /etc/openvpn/secret.key"); err != nil {
		log.Fatal(err)
	}

	return b.String()
}

func waitForSetup(ip string) *ssh.Client {
	for {
		client, err := createSSHClient(ip)
		if err == nil {
			session, err := client.NewSession()
			if err == nil {
				defer session.Close()
				err = session.Run("cat /root/yovpn.ready")
				if err == nil {
					return client
				}
			}
		}
		log.Printf("Waiting for SSH connection... (%s)\n", err)
		<-time.After(time.Second * 10)
	}
}
