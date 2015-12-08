// Package web contains the HTTP endpoints used by the yovpn-server binary.
package web

import (
	"net/http"

	"github.com/xperimental/yovpn/provisioner"
)

// SetupHandlers sets up the handler functions for a provisioner.
func SetupHandlers(provisioner provisioner.Provisioner) {
	http.HandleFunc("/cleanup", cleanupHandler(provisioner))
	http.HandleFunc("/endpoint", endpointHandler(provisioner))
	http.HandleFunc("/provision", startProvision(provisioner))
	http.HandleFunc("/regions", regionsHandler(provisioner))
	http.HandleFunc("/", blankPage)
}
