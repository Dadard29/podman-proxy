package podman

import (
	"fmt"

	"github.com/Dadard29/podman-proxy/models"
)

// ListContainers retrieve all containers
func (r *PodmanRuntime) ListContainers() ([]*models.PodmanContainer, error) {
	containers, err := r.runtime.GetAllContainers()
	if err != nil {
		return nil, err
	}
	out := make([]*models.PodmanContainer, 0)
	for _, container := range containers {

		podmanContainer, err := models.NewPodmanContainer(container)
		if err != nil {
			return nil, err
		}

		out = append(out, podmanContainer)
	}

	return out, nil
}

// GetContainerFromName retrieve a specific container using its name
func (r *PodmanRuntime) GetContainerFromName(containerName string) (*models.PodmanContainer, error) {
	var out *models.PodmanContainer

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
func (r *PodmanRuntime) GetContainerFromID(containerId string) (*models.PodmanContainer, error) {
	container, err := r.runtime.GetContainer(containerId)
	if err != nil {
		return nil, err
	}

	podmanContainer, err := models.NewPodmanContainer(container)
	if err != nil {
		return nil, err
	}

	return podmanContainer, nil
}
