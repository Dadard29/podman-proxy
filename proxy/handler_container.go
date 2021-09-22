package proxy

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// List all containersPodman stored in database
func (p *Proxy) containersGet(w http.ResponseWriter, r *http.Request) {
	containersPodman, err := p.db.ListContainers()
	if err != nil {
		p.logger.Println(err)
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	p.WriteJson(w, &containersPodman)
}

// Main entrypoint for containers list management
func (p *Proxy) containersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.containersGet(w, r)
	}
}

// Get a specific container from name
func (p *Proxy) containerGet(w http.ResponseWriter, r *http.Request, containerName string) {
	container, err := p.db.GetContainer(containerName)
	if err != nil {
		p.logger.Println(err)
		p.WriteErrorJson(w, http.StatusNotFound, err)
		return
	}

	p.WriteJson(w, &container)
}

// Set the exposed port by an existing container
func (p *Proxy) containerPost(w http.ResponseWriter, r *http.Request, containerName string) {
	exposedPortStr := r.URL.Query().Get("exposedPort")
	exposedPort, err := strconv.Atoi(exposedPortStr)
	if err != nil {
		p.WriteErrorJson(w, http.StatusBadRequest, err)
		return
	}

	err = p.db.UpdateContainerExposedPort(containerName, exposedPort)
	if err != nil {
		p.WriteErrorJson(w, http.StatusNotFound, err)
		return
	}

	container, err := p.db.GetContainer(containerName)
	if err != nil {
		p.WriteErrorJson(w, http.StatusNotFound, err)
		return
	}

	p.WriteJson(w, &container)
}

// Main entrypoint for single container management
func (p *Proxy) containerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerName := vars["container"]

	if r.Method == http.MethodGet {
		p.containerGet(w, r, containerName)

	} else if r.Method == http.MethodPost {
		p.containerPost(w, r, containerName)

	}
}
