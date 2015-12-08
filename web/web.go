// Package web contains the HTTP endpoints used by the yovpn-server binary.
package web

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/xperimental/yovpn/provisioner"
)

// CreateServer sets up a HTTP server for a provisioner.
func CreateServer(provisioner provisioner.Provisioner) http.Handler {
	e := echo.New()

	e.Get("/cleanup", cleanupHandler(provisioner))
	e.Get("/endpoint", endpointHandler(provisioner))
	e.Get("/provision", startProvision(provisioner))
	e.Get("/regions", regionsHandler(provisioner))
	e.Get("/", blankPage)

	return e
}
