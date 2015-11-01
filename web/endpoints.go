package web

import (
	"net/http"

	"github.com/xperimental/yovpn/provisioner"
)

func EndpointHandler(provisioner *provisioner.Provisioner) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if len(id) == 0 {
			writeJSON(w, 400, jsonError{Error: "ID not provided!"})
			return
		}

		endpoint, err := provisioner.GetEndpoint(id)
		if err != nil {
			writeJSON(w, 404, jsonError{Error: "Endpoint not found!"})
			return
		}

		writeJSON(w, 200, endpoint)
	}
}

func CleanupHandler(provisioner *provisioner.Provisioner) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoints := provisioner.ListEndpoints()
		for _, endpoint := range endpoints {
			provisioner.DestroyEndpoint(endpoint.ID)
		}
	}
}
