package main

import (
	"log"
	"net/http"

	"github.com/xperimental/yovpn/config"
	"github.com/xperimental/yovpn/provisioner"
	"github.com/xperimental/yovpn/web"
)

func main() {
	config := config.GetConfig()

	log.Println("Create provisioner...")
	provisioner, err := provisioner.NewProvisioner(config.Token)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Setup handlers...")
	web.SetupHandlers(provisioner)
	http.HandleFunc("/", web.BlankPage)

	log.Printf("Listen on %s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
