package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

var keyfile = flag.String("keyfile", "yovpn.key", "SSH keyfile to use for connection")

func readKey() (key []byte, err error) {
	file, err := os.Open(*keyfile)
	switch {
	case os.IsNotExist(err):
		err = errKeyNotFound
		return
	case err != nil:
		return
	}
	defer file.Close()

	key, err = ioutil.ReadAll(file)
	return
}

func loadPrivateKey() (ssh.Signer, error) {
	key, err := readKey()
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(key)
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

func createPrivateKey() ssh.Signer {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPem := &pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}
	pemBytes := pem.EncodeToMemory(privateKeyPem)
	err = ioutil.WriteFile(*keyfile, pemBytes, 0600)
	if err != nil {
		log.Fatal(err)
	}

	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		log.Fatal(err)
	}
	return signer
}
