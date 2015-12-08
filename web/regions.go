package web

import (
	"net/http"

	"github.com/xperimental/yovpn/provisioner"
)

func regionsHandler(provisioner provisioner.Provisioner) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		regions, err := provisioner.ListRegions()
		if err != nil {
			handleError(w, err, "Could not list regions!")
			return
		}

		writeJSON(w, 200, regions)
	}
}
