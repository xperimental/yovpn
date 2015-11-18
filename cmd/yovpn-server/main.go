package main

import (
	"net/http"

	"github.com/prometheus/common/log"
	"github.com/xperimental/yovpn/config"
	"github.com/xperimental/yovpn/provisioner"
	"github.com/xperimental/yovpn/web"
)

func backgroundRunner(p *provisioner.Provisioner) {
	select {
	case <-p.Signal:
		log.Debugf("Signal from provisioner.")
	}
}

func main() {
	config := config.GetConfig()

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
