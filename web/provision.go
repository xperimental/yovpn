package web

import (
	"net/http"

	"github.com/xperimental/yovpn/provisioner"
)

func StartProvision(provisioner *provisioner.Provisioner) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		region := r.URL.Query().Get("region")
		if len(region) == 0 {
			w.WriteHeader(400)
			w.Write([]byte("You need to provide a region!"))
			return
		}

		endpoint, err := provisioner.CreateEndpoint(region)
		if err != nil {
			handleError(w, err, "Failed to create endpoint!")
			return
		}

		writeJSON(w, 201, endpoint)
	}
}
