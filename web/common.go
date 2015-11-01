package web

import (
	"encoding/json"
	"log"
	"net/http"
)

func handleError(w http.ResponseWriter, err error, msg string) {
	log.Printf("ERROR \"%s\": %+v", msg, err)

	w.WriteHeader(500)
	w.Write([]byte(msg))
}

func writeJSON(w http.ResponseWriter, status int, content interface{}) {
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		handleError(w, err, "Failed to create JSON output!")
		return
	}

	w.WriteHeader(201)
	w.Write(data)
}
