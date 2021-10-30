package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Dadard29/podman-proxy/proxy"
)

func serve(p *proxy.Proxy) {

	err := p.Serve(false)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(Fatal(err))
	}
}

func main() {
	p, err := proxy.NewProxy()
	if err != nil {
		log.Fatal(Fatal(err))
	}

	// log.Println(Info("* Starting upgrader goroutine"))
	// go func() {
	// 	err := p.Upgrader.Serve()
	// 	if err != nil && err != http.ErrServerClosed {
	// 		log.Fatal(Fatal(err))
	// 	}
	// }()

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
