package main

import (
	"github.com/Dadard29/podman-proxy/proxy"
)

func main() {
	conf := proxy.RetrieveConfig()
	p := proxy.NewProxy(conf)

	// blocking method (run the proxy http server)
	p.Start()
}

