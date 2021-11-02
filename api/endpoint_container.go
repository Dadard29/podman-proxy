package api

import (
	"net/http"
	"strconv"

	"github.com/Dadard29/podman-proxy/models"
)

// List all containersPodman stored in database
func (a *Api) ContainersGet(w http.ResponseWriter, r *http.Request) (*[]models.Container, error) {
	containersPodman, err := a.db.ListContainers()
	if err != nil {
		return nil, err
	}

	return &containersPodman, nil
}

// Get a specific container from name
func (a *Api) ContainerGet(w http.ResponseWriter, r *http.Request, containerName string) (*models.Container, error) {
	container, err := a.db.GetContainer(containerName)
	if err != nil {
		return nil, err
	}

	return &container, nil
}

// Set the exposed port by an existing container
func (a *Api) ContainerPost(w http.ResponseWriter, r *http.Request, containerName string) (*models.Container, error) {
	exposedPortStr := r.URL.Query().Get("exposedPort")
	exposedPort, err := strconv.Atoi(exposedPortStr)
	if err != nil {
		return nil, err
	}

	err = a.db.UpdateContainerExposedPort(containerName, exposedPort)
	if err != nil {
		return nil, err
	}

	container, err := a.db.GetContainer(containerName)
	if err != nil {
		return nil, err
	}

	return &container, nil
}
