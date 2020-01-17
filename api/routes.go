package api

import "net/http"

type Route struct {
	HttpMethods []string
	Handler func(w http.ResponseWriter, r *http.Request)
}

var routes = map[string]Route{
	"/rules": {
		HttpMethods: []string{http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodPut},
		Handler:     rulesHandler,
	},
}

func GetRoutes() map[string]Route {
	return routes
}
