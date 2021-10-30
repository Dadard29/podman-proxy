package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Dadard29/podman-proxy/proxy"
)

const InfraLogFreq = 5 * time.Minute

func serve(p *proxy.Proxy) {

	err := p.Serve(true)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(Fatal(err))
	}
}

func upgraderServe(p *proxy.Proxy) {
	err := p.Upgrader.Serve()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(Fatal(err))
	}
}

func infraLogger(p *proxy.Proxy) {
	for {
		if err := p.NewInfraLog(); err != nil {
			log.Println(Warn("failed inserting to infra log", err))
		}
		time.Sleep(InfraLogFreq)
	}
}

func main() {
	p, err := proxy.NewProxy()
	if err != nil {
		log.Fatal(Fatal(err))
	}

	log.Println(Info("* Starting infra logger goroutine"))
	go infraLogger(p)

	log.Println(Info("* Starting upgrader goroutine"))
	go upgraderServe(p)

	log.Println(Info("* Starting proxy goroutine"))
	go serve(p)

	defer func() {
		log.Println(Info("* Shutting down.."))
		// p.Upgrader.Shutdown()
		p.Shutdown()
		p.Close()
	}()

	// running routines until interrupt is received
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for {
		for p.Ctx == nil {
			continue
		}

		select {
		case <-(*p.Ctx).Done():
			log.Println(Info("* Restarting.."))
			p.Shutdown()
			go serve(p)
		case <-c:
			return
		}
	}
}
