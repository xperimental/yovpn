package provisioner

import (
	"bytes"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

func createSSHClient(signer ssh.Signer, ip string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	return ssh.Dial("tcp", net.JoinHostPort(ip, "22"), config)
}

func readSecret(client *ssh.Client) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("cat /etc/openvpn/secret.key"); err != nil {
		return "", err
	}

	return b.String(), nil
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
