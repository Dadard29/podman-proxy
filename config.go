package main

import (
	"log"
	"os"
)

const (
	proxyHostKey = "PODMAN_PROXY_HOST"
	proxyPortKey = "PODMAN_PROXY_PORT"
)

type config struct {
	// the host which will not be redirected to podman container
	// the api will be available through this host
	ProxyHost string `json:"proxy_host"`
	ProxyPort int `json:"proxy_port"`
}

func getDefaultConfig() config {
	defaultProxyHost, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}

	defaultProxyPort := 8080

	return config{
		ProxyHost: defaultProxyHost,
		ProxyPort: defaultProxyPort,
	}
}

func retrieveEnv() config {
	c := getDefaultConfig()

	if proxyHost := os.Getenv(proxyHostKey); proxyHost != "" {
		c.ProxyHost = proxyHost
	}

	//if proxyPort := os.Getenv(proxyPortKey); proxyPort != "" {
	//	c.ProxyPort = proxyPort
	//}

	return c
}
