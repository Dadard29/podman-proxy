package main

import (
	"context"
	"fmt"
	"github.com/containers/libpod/libpod"
)

func testApiPodman() {
	c := context.Background()
	runtime, err := libpod.NewRuntime(c)
	if err != nil {
		fmt.Println("error setup")
		fmt.Println(err)
		return
	}

	containers, err := runtime.GetAllContainers()
	if err != nil {
		fmt.Println("error getting all containers")
		fmt.Println(err)
		return
	}

	if containers != nil {
		for _, con := range containers {
			fmt.Println(con.Config().Spec.Process.Env)
		}
	}

}

func main() {
	conf := retrieveEnv()
	startProxy(conf)
}

