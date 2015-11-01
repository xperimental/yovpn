package provisioner

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func publicFingerprint(key ssh.Signer) string {
	h := md5.New()
	h.Write(key.PublicKey().Marshal())
	sum := h.Sum(nil)
	var buf bytes.Buffer
	for i, b := range sum {
		bs := fmt.Sprintf("%x", b)
		if len(bs) == 1 {
			buf.WriteRune('0')
		}
		buf.WriteString(bs)
		if i < len(sum)-1 {
			buf.WriteRune(':')
		}
	}
	return buf.String()
}

func createPrivateKey() (signer ssh.Signer, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}

	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPem := &pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}

	pemBytes := pem.EncodeToMemory(privateKeyPem)
	signer, err = ssh.ParsePrivateKey(pemBytes)
	return
}
