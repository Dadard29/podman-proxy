package main

import (
	"fmt"
	"github.com/Dadard29/podman-proxy/api"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var globalConf config

// this handler will redirect the requests to the container
// this is the proxy
func mainProxyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte ("main proxy handler\n"))
}

func startProxy(conf config) {
	globalConf = conf
	proxyUrl := fmt.Sprintf(":%d", globalConf.ProxyPort)

	router := mux.NewRouter()


	// setup the main proxy route
	router.HandleFunc("/", mainProxyHandler)

	// setup the proxy web api
	for p, r := range api.GetRoutes() {
		router.HandleFunc(p, r.Handler).Methods(r.HttpMethods...)
	}


	log.Printf("listening on %s...\n", proxyUrl)
	server := &http.Server{
		Addr:              proxyUrl,
		Handler:           router,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          log.New(os.Stdout, "podman-proxy.db", log.Ldate | log.Ltime),
		BaseContext:       nil,
		ConnContext:       nil,
	}
	err := server.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}
}
