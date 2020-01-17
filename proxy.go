package main

import (
	"fmt"
	"github.com/Dadard29/podman-proxy/api"
	"log"
	"net/http"
	"os"
)

var globalConf config

func checkMethod(methods []string, requestMethod string) bool {
	for _, v := range methods {
		if v == requestMethod {
			return true
		}
	}
	return false
}

// this handler will redirect the requests to the container
// this is the proxy
func mainProxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Host == globalConf.getAddr() {

		// redirect the request to the correct handler
		for u, h := range api.GetRoutes() {
			if u == r.URL.String() && checkMethod(h.HttpMethods, r.Method) {
				h.Handler(w, r)
				return
			}
		}

		// no valid handler found for this route and this http method
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte (fmt.Sprintf("%d %s\n", http.StatusNotFound, "custom page not found")))
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		// retrieve the container ip and port from the request Host
		_, err := w.Write([]byte ("redirected to some container\n"))
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func startProxy(conf config) {
	globalConf = conf
	proxyUrl := fmt.Sprintf(":%d", globalConf.ProxyPort)

	// setup the main proxy route
	http.HandleFunc("/", mainProxyHandler)

	log.Printf("listening on %s...\n", proxyUrl)
	server := &http.Server{
		Addr:              proxyUrl,
		Handler:           nil,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          log.New(os.Stdout, "podman-proxy ", log.Ldate | log.Ltime),
		BaseContext:       nil,
		ConnContext:       nil,
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}
}
