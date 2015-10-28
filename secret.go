package main

import (
	"bytes"
	"crypto/md5"
	"flag"
	"fmt"
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

func loadPrivateKey() ssh.Signer {
	key := readKey()
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal(err)
	}
	return signer
}

func publicFingerprint(key ssh.Signer) string {
	h := md5.New()
	h.Write(key.PublicKey().Marshal())
	sum := h.Sum(nil)
	var buf bytes.Buffer
	for i, b := range sum {
		buf.WriteString(fmt.Sprintf("%x", b))
		if i < len(sum)-1 {
			buf.WriteRune(':')
		}
	}
	return buf.String()
}

func createSSHClient(signer ssh.Signer, ip string) (*ssh.Client, error) {
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

func waitForSetup(signer ssh.Signer, ip string) *ssh.Client {
	for {
		client, err := createSSHClient(signer, ip)
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
