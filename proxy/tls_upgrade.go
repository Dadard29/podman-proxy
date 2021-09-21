package proxy

import (
	"log"
	"net/http"
	"os"
)

func (p *Proxy) upgrader(w http.ResponseWriter, req *http.Request) {
	// https://gist.github.com/d-schmidt/587ceec34ce1334a5e60

	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	p.logger.Printf("redirecting to %s...\n", target)
	http.Redirect(w, req, target,
		http.StatusMovedPermanently)
}

func (p *Proxy) UpgraderServe() error {
	// fixme
	server := http.Server{
		Addr:     ":9001",
		Handler:  http.HandlerFunc(p.upgrader),
		ErrorLog: log.New(os.Stdout, "podman-proxy-upgrader ", log.Ldate|log.Ltime),
	}
	p.logger.Printf("Starting proxy upgrader on %s\n", server.Addr)
	return server.ListenAndServe()
}
