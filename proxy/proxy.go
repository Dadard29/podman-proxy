package proxy

import (
	"fmt"
	"github.com/Dadard29/podman-proxy/api"
	"log"
	"net/http"
	"os"
)

var globalProxy *Proxy

type Proxy struct {
	exposedApi *api.Api
	httpServer *http.Server
	config     Config
}

// this handler will redirect the requests to the container
// this is the proxy
func mainProxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Host == globalProxy.config.getAddr() {

		// redirect the request to the correct handler
		for u, h := range globalProxy.exposedApi.GetRoutes() {
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

func NewProxy(conf Config) *Proxy {
	proxyUrl := fmt.Sprintf("%s:%d", conf.ProxyHost, conf.ProxyPort)

	// setup the main proxy route
	http.HandleFunc("/", mainProxyHandler)

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
		ErrorLog:          log.New(os.Stdout, "podman-proxy ", log.Ldate|log.Ltime),
		BaseContext:       nil,
		ConnContext:       nil,
	}

	exposeApi := api.NewApi()

	p := &Proxy{
		exposedApi: exposeApi,
		httpServer: server,
		config: conf,
	}

	globalProxy = p
	return p
}

func (p *Proxy) Start() {
	log.Printf("listening on %s...\n", p.httpServer.Addr)
	err := p.httpServer.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}
}


