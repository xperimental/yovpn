package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"

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

func readSecret(ip string) string {
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
	client, err := ssh.Dial("tcp", net.JoinHostPort(ip, "22"), config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

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
