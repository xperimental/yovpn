package web

import "net/http"

func BlankPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
