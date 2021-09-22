package podman

import (
	"fmt"

	"github.com/Dadard29/podman-proxy/models"
)

// ListContainers retrieve all containers
func (r *PodmanRuntime) ListContainers() ([]*models.Container, error) {
	containers, err := r.runtime.GetAllContainers()
	if err != nil {
		return nil, err
	}
	out := make([]*models.Container, 0)
	for _, container := range containers {

		Container, err := models.NewContainer(container)
		if err != nil {
			return nil, err
		}

		out = append(out, Container)
	}

	return out, nil
}

// GetContainerFromName retrieve a specific container using its name
func (r *PodmanRuntime) GetContainerFromName(containerName string) (*models.Container, error) {
	var out *models.Container

	containers, err := r.ListContainers()
	if err != nil {
		return out, err
	}

	for _, container := range containers {
		if container.Name == containerName {
			return container, nil
		}
	}

	return out, fmt.Errorf("container with name %s not found", containerName)
}

// GetContainerFromID retrieve a specific container using its ID
func (r *PodmanRuntime) GetContainerFromID(containerId string) (*models.Container, error) {
	container, err := r.runtime.GetContainer(containerId)
	if err != nil {
		return nil, err
	}

	Container, err := models.NewContainer(container)
	if err != nil {
		return nil, err
	}

	return Container, nil
}
