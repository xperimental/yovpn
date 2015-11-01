package web

import (
	"net/http"

	"github.com/xperimental/yovpn/provisioner"
)

func StartProvision(provisioner *provisioner.Provisioner) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		region := r.URL.Query().Get("region")
		if len(region) == 0 {
			writeJSON(w, 400, jsonError{Error: "You need to provide a region!"})
			return
		}

		endpoint := provisioner.CreateEndpoint(region)
		writeJSON(w, 201, endpoint)
	}
}
