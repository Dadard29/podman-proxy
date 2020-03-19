package api

import (
	"github.com/containers/libpod/libpod"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

var globalApi *Api

type Api struct {
	connector     *gorm.DB
	libpodRuntime *libpod.Runtime
	routes        map[string]Route
}

// api object to be used by the api handlers
// methods are supposed to used only by the api handlers
// static function are supposed to be used by other packages

func NewApi() *Api {
	con := newConnector()

	runtime := NewLibpodRuntime()

	routes := map[string]Route{
		"/rules": {
			HttpMethods: []string{http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodPut},
			Handler:     rulesHandler,
		},
		"/rules/list": {
			HttpMethods: []string{http.MethodGet},
			Handler:     rulesListHandler,
		},
	}

	a := &Api{
		connector:     con,
		libpodRuntime: runtime,
		routes:        routes,
	}

	globalApi = a
	return a
}

func (a *Api) Close() {
	err := a.connector.Close()
	if err != nil {
		log.Fatalln(err)
	}
}

func (a *Api) GetRoutes() map[string]Route {
	return a.routes
}

func (a *Api) PingDb() bool {
	return a.connector.DB().Ping() == nil
}
