package proxy

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Update the database with all existing containersPodman
func (p *Proxy) containersPut(w http.ResponseWriter, r *http.Request) {
	containersDb, err := p.db.ListContainers()
	if err != nil {
		p.logger.Println(err)
		p.WriteErrorJson(w, err)
		return
	}

	containersPodman, err := p.podman.ListContainers()
	if err != nil {
		p.logger.Println(err)
		p.WriteErrorJson(w, err)
		return
	}

	for _, containerPodman := range containersPodman {
		found := false

		// updating the database with the new IPs of the containers
		for _, containerDb := range containersDb {
			if containerPodman.Name == containerDb.Name {
				err := p.db.UpdateContainer(containerPodman.Name, containerPodman.IpAddress, containerPodman.Status)
				if err != nil {
					p.logger.Println(err)
				}
				found = true
				break
			}
		}

		// creating in database the newly created containers
		if !found {
			err := p.db.InsertContainer(containerPodman)
			if err != nil {
				p.logger.Println(err)
			}
		}
	}

	// deleting the non-existant containers in database
	for _, containerDb := range containersDb {
		found := false

		for _, containerPodman := range containersPodman {
			if containerPodman.Name == containerDb.Name {
				found = true
			}
		}

		if !found {
			err := p.db.DeleteContainer(containerDb.Name)
			if err != nil {
				p.logger.Println(err)
			}
		}
	}

	// retrieve the new database containers
	containersDb, err = p.db.ListContainers()
	if err != nil {
		if err != nil {
			p.logger.Println(err)
			p.WriteErrorJson(w, err)
			return
		}
	}

	p.WriteJson(w, &containersDb)

}

// List all containersPodman stored in database
func (p *Proxy) containersGet(w http.ResponseWriter, r *http.Request) {
	containersPodman, err := p.db.ListContainers()
	if err != nil {
		p.logger.Println(err)
		p.WriteErrorJson(w, err)
		return
	}

	p.WriteJson(w, &containersPodman)
}

// Main entrypoint for containers list management
func (p *Proxy) containersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.containersGet(w, r)

	} else if r.Method == http.MethodPut {
		p.containersPut(w, r)

	}
}

// Get a specific container from name
func (p *Proxy) containerGet(w http.ResponseWriter, r *http.Request, containerName string) {
	container, err := p.db.GetContainer(containerName)
	if err != nil {
		p.logger.Println(err)
		p.WriteErrorJson(w, err)
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
	}
}
