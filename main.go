package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Dadard29/podman-proxy/proxy"
)

func main() {
	p, err := proxy.NewProxy()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("* Starting upgrader goroutine")
	go func() {
		err := p.UpgraderServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Println("* Starting proxy goroutine")
	go func() {
		err := p.Serve(false)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	defer func() {
		log.Println("* Shutting down proxy..")
		p.Shutdown()
	}()

	// running routines until interrupt is received
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
