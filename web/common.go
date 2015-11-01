package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/xperimental/yovpn/provisioner"
)

type jsonError struct {
	Error string `json:"error"`
}

func handleError(w http.ResponseWriter, err error, msg string) {
	log.Printf("ERROR \"%s\": %+v", msg, err)

	writeJSON(w, 500, jsonError{Error: msg})
}

func writeJSON(w http.ResponseWriter, status int, content interface{}) {
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		handleError(w, err, "Failed to create JSON output!")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)
}

func SetupHandlers(provisioner *provisioner.Provisioner) {
	http.HandleFunc("/cleanup", CleanupHandler(provisioner))
	http.HandleFunc("/endpoint", EndpointHandler(provisioner))
	http.HandleFunc("/provision", StartProvision(provisioner))
	http.HandleFunc("/regions", RegionsHandler(provisioner))
}
