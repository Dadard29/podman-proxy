package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/containers/libpod/libpod"
	"log"
)

func NewLibpodRuntime() *libpod.Runtime {
	c := context.Background()
	runtime, err := libpod.NewRuntime(c)
	if err != nil {
		log.Fatalln(err)
	}

	return runtime
}

func (a *Api) GetContainerFromLibpod(containerName string) (*libpod.Container, error) {

	containers, err := a.libpodRuntime.GetAllContainers()
	if err != nil {
		return nil, err
	}

	if containers != nil {
		for _, con := range containers {
			if con.Name() == containerName {
				return con, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("the container with name %s does not exists", containerName))
}

func (a *Api) GetContainerIp(con *libpod.Container) (string, error) {
	ips, err := con.IPs()
	if err != nil {
		return "", err
	}

	return ips[0].IP.String(), nil
}
