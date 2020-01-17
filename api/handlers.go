package api

import "net/http"

func rulesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte ("rule handler\n"))
}
