package podman

import "github.com/Dadard29/podman-proxy/models"

func (r *PodmanRuntime) ListContainers() ([]models.Container, error) {
	containers, err := r.runtime.GetAllContainers()
	if err != nil {
		return nil, err
	}
	out := make([]models.Container, 0)
	for _, container := range containers {
		name := container.Name()
		ips, err := container.IPs()
		if err != nil {
			return nil, err
		}
		ip := ips[0].IP.String()

		ports, err := container.PortMappings()
		if err != nil {
			return nil, err
		}
		exposedPort := int(ports[0].ContainerPort)

		out = append(out, models.Container{
			Name:        name,
			IpAddress:   ip,
			ExposedPort: exposedPort,
		})
	}

	return out, nil
}
