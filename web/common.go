package web

import (
	"encoding/json"
	"log"
	"net/http"
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
