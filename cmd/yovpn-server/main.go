// This package provides a binary that provides a HTTP endpoint to be used for provisioning multiple VPN endpoints.
package main

import (
	"net/http"

	"github.com/prometheus/common/log"
	"github.com/xperimental/yovpn/provisioner"
	"github.com/xperimental/yovpn/web"
)

func backgroundRunner(p provisioner.Provisioner) {
	select {
	case <-p.Signal():
		log.Debugf("Signal from provisioner.")
	}
}

func main() {
	config := createConfig()

	log.Info("Create provisioner...")
	provisioner, err := provisioner.NewProvisioner(config.Token)
	if err != nil {
		log.Fatal(err)
	}
	go backgroundRunner(provisioner)

	log.Info("Setup handlers...")
	web.SetupHandlers(provisioner)
	http.HandleFunc("/", web.BlankPage)

	log.Infof("Listen on %s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
