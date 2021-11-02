package proxy

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (p *Proxy) containerGetHandler(w http.ResponseWriter, r *http.Request, containerName string) {
	container, err := p.api.ContainerGet(w, r, containerName)
	if err != nil {
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	p.WriteJson(w, container)
}

func (p *Proxy) containerPostHandler(w http.ResponseWriter, r *http.Request, containerName string) {
	container, err := p.api.ContainerPost(w, r, containerName)
	if err != nil {
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	p.WriteJson(w, container)
}

// Main entrypoint for single container management
func (p *Proxy) containerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerName := vars["container"]

	if r.Method == http.MethodGet {
		p.containerGetHandler(w, r, containerName)

	} else if r.Method == http.MethodPost {
		p.containerPostHandler(w, r, containerName)

	}
}
