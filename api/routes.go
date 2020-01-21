package api

import "net/http"

type Route struct {
	HttpMethods []string
	Handler func(w http.ResponseWriter, r *http.Request)
}
