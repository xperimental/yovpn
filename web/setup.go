package web

import (
	"net/http"

	"github.com/xperimental/yovpn/provisioner"
)

func SetupHandlers(provisioner provisioner.Provisioner) {
	http.HandleFunc("/cleanup", CleanupHandler(provisioner))
	http.HandleFunc("/endpoint", EndpointHandler(provisioner))
	http.HandleFunc("/provision", StartProvision(provisioner))
	http.HandleFunc("/regions", RegionsHandler(provisioner))
}
