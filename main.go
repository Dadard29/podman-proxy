package main

import (
	"context"
	"fmt"
	"github.com/containers/libpod/cmd/podman/shared"
	"github.com/containers/libpod/libpod"
	createconfig "github.com/containers/libpod/pkg/spec"
)

func main() {
	c := context.Background()
	runtime, err := libpod.NewRuntime(c)
	if err != nil {
		fmt.Println("error setup")
		fmt.Println(err)
		return
	}

	containers, err := runtime.GetAllContainers()
	if err != nil {
		fmt.Println("error get all containers")
		fmt.Println(err)
		return
	}

	if containers != nil {
		for _, con := range containers {
			fmt.Println(con.Config().Spec.Process.Env)
		}
	}


	cc := createconfig.CreateConfig{}
	cc.Image = "nginx"
	cc.Name = "new_server"

	var pod *libpod.Pod
	container, err := shared.CreateContainerFromCreateConfig(runtime, &cc, context.Background(), pod)
	if err != nil {
		fmt.Println("error creating container")
		fmt.Println(err)
		return
	}

	fmt.Println(container.Name(), container.ID())

}

