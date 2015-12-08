package web

import "net/http"

func blankPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
