package proxy

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Upgrader struct {
	port   int
	server *http.Server
	logger *log.Logger
}

func (up *Upgrader) upgraderHandler(w http.ResponseWriter, req *http.Request) {
	// https://gist.github.com/d-schmidt/587ceec34ce1334a5e60

	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	http.Redirect(w, req, target, http.StatusMovedPermanently)
}

func (up *Upgrader) upgraderMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		up.logger.Printf("%s %s %s", r.Method, r.Host, r.URL.Path)
	}

	return http.HandlerFunc(fn)
}

func NewUpgrader(port int) *Upgrader {

	addr := fmt.Sprintf(":%d", port)

	server := http.Server{
		Addr:     addr,
		Handler:  nil,
		ErrorLog: log.New(os.Stdout, "upgrader-error ", log.Ldate|log.Ltime),
	}

	logger := log.New(log.Default().Writer(), "upgrader ", log.Default().Flags())

	up := &Upgrader{
		port:   port,
		server: &server,
		logger: logger,
	}

	router := mux.NewRouter()
	router.PathPrefix("/").HandlerFunc(up.upgraderHandler)
	router.Use(up.upgraderMiddleware)
	up.server.Handler = router

	return up
}

func (up *Upgrader) Serve() error {

	up.logger.Printf("Starting proxy upgrader on %s\n", up.server.Addr)
	return up.server.ListenAndServe()
}
