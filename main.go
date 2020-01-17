package main

import (
	"context"
	"fmt"
	"github.com/containers/libpod/libpod"
	"net/http"
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

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("requested")
	w.Write([]byte ("salut"))
}

func main() {
	conf := retrieveEnv()
	startProxy(conf)

	http.HandleFunc("/", proxyHandler)
	http.ListenAndServe(":8080", nil)
}

