package models

import (
	"database/sql"
	"fmt"

	"github.com/containers/libpod/libpod"
)

type Container struct {
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

	// the port exposed by the container
	// set manually by the user
	ExposedPort int `json:"exposed_port"`

	// the container status
	Status ContainerStatus `json:"state"`
}

func NewContainer(container *libpod.Container) (*Container, error) {
	isInPod := container.PodID() != ""
	isInfra := container.IsInfra()

	state, err := container.ContainerState()
	if err != nil {
		return nil, err
	}
	statusStr := state.State.String()
	status := NewContainerStatus(statusStr)

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

	return &Container{
		Id:          container.ID(),
		Name:        container.Name(),
		IsInfra:     isInfra,
		IsInPod:     isInPod,
		PodId:       container.PodID(),
		IpAddress:   ip,
		ExposedPort: 0,
		Status:      status,
	}, nil
}

func (c Container) String() string {
	return fmt.Sprintf("%s: %s:%d", c.Name, c.IpAddress, c.ExposedPort)
}

func NewContainerFromRow(rowCursor *sql.Rows) (Container, error) {
	var id string
	var name string
	var isInfra int
	var isInPod int
	var podId string
	var ipAddress string
	var exposedPort int
	var status string
	err := rowCursor.Scan(&id, &name, &isInfra, &isInPod, &podId,
		&ipAddress, &exposedPort, &status)

	if err != nil {
		return Container{}, err
	}

	return Container{
		Id:          id,
		Name:        name,
		IsInfra:     isInfra == 1,
		IsInPod:     isInPod == 1,
		PodId:       podId,
		IpAddress:   ipAddress,
		ExposedPort: exposedPort,
		Status:      NewContainerStatus(status),
	}, nil
}
