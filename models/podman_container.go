package models

import (
	"fmt"

	"github.com/containers/libpod/libpod"
)

type PodmanContainer struct {
	// the container ID
	Id string `json:"id"`

	// the container name
	Name string `json:"name"`

	// true if it is an pod infra container
	IsInfra bool `json:"is_infra"`

	// true if belongs to a pod, false else
	IsInPod bool `json:"is_in_pod"`

	// ID of the pod it belongs to. "" if standalone
	PodId string `json:"pod_id"`

	// the container IP address in the podman network
	IpAddress string `json:"ip_address"`

	// the container status
	Status PodmanContainerStatus `json:"state"`
}

func NewPodmanContainer(container *libpod.Container) (*PodmanContainer, error) {
	isInPod := container.PodID() != ""
	isInfra := container.IsInfra()

	state, err := container.ContainerState()
	if err != nil {
		return nil, err
	}
	statusStr := state.State.String()
	status := NewPodmanContainerStatus(statusStr)

	var ip string
	// checking pod context
	// able to get IP only if out of a pod, or is an infra
	if !isInPod || (isInPod && isInfra) {

		// checking status
		// able to get IP only on running containers
		if status == Running {

			ips, err := container.IPs()
			if err != nil {
				return nil, err
			}
			if len(ips) == 0 {
				// no IP detected for this container
				// container might be running rootless or in incompatible state
				return nil, fmt.Errorf("no IP detected for container with name %s", container.Name())
			}
			ip = ips[0].IP.String()
		}
	}

	return &PodmanContainer{
		Id:        container.ID(),
		Name:      container.Name(),
		IsInfra:   isInfra,
		IsInPod:   isInPod,
		PodId:     container.PodID(),
		IpAddress: ip,
		Status:    status,
	}, nil
}
