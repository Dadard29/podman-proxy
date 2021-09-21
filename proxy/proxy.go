package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Dadard29/podman-proxy/db"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
)

type Proxy struct {
	config config
	db     *db.Db
	logger *log.Logger
	server *http.Server
	router *mux.Router
}

func (p *Proxy) getAddrFromProxyPort() string {
	return fmt.Sprintf(":%d", p.config.proxyPort)
}

func (p *Proxy) Host() string {
	if p.config.proxyPort == 80 || p.config.proxyPort == 443 {
		return p.config.proxyHost
	} else {
		return fmt.Sprintf("%s:%d", p.config.proxyHost, p.config.proxyPort)
	}
}

func (p *Proxy) newHttpsServer(domainNames ...string) *http.Server {
	// tls
	manager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("/srv/https/certificates"),
		HostPolicy: autocert.HostWhitelist(domainNames...),
	}

	// server
	addr := p.getAddrFromProxyPort()
	server := &http.Server{
		Addr:      addr,
		TLSConfig: manager.TLSConfig(),
		ErrorLog:  log.New(os.Stdout, "podman-proxy-error ", log.Ldate|log.Ltime),
	}

	return server
}

func (p *Proxy) newHttpServer() *http.Server {

	// server
	addr := p.getAddrFromProxyPort()
	return &http.Server{
		Addr:     addr,
		ErrorLog: log.New(os.Stdout, "podman-proxy-error ", log.Ldate|log.Ltime),
	}
}

func NewProxy() (*Proxy, error) {
	config, err := newConfigFromEnv()
	if err != nil {
		return nil, err
	}

	proxyDb, err := db.NewDb(config.dbPath)
	if err != nil {
		return nil, err
	}

	logger := log.New(log.Default().Writer(), "podman-proxy ", log.Default().Flags())

	proxy := &Proxy{
		config: config,
		db:     proxyDb,
		logger: logger,
		server: nil,
		router: nil,
	}

	router := mux.NewRouter()
	router.HandleFunc("/rule", proxy.rulesHandler).Methods(http.MethodGet)
	router.HandleFunc("/rule/{dn}", proxy.ruleHandler).Methods(http.MethodGet, http.MethodPost, http.MethodDelete)
	router.HandleFunc("/domain-name", proxy.domainNamesHandler).Methods(http.MethodGet)
	router.HandleFunc("/domain-name/{dn}", proxy.domainNameHandler).Methods(http.MethodGet, http.MethodPost, http.MethodDelete)
	router.HandleFunc("/container", proxy.containersHandler).Methods(http.MethodGet)
	router.HandleFunc("/container/{container}", proxy.containerHandler).Methods(http.MethodGet, http.MethodPost, http.MethodDelete)

	proxy.router = router

	return proxy, nil
}

func (p *Proxy) Serve(withTLS bool) error {

	router := mux.NewRouter()
	router.PathPrefix("/").HandlerFunc(p.switcher)
	router.Use(p.dbLoggingMiddleware)

	p.logger.Printf("Starting proxy server on %s...\n", p.Host())

	if withTLS {
		// get the list of domain names to register
		domainNamesList, err := p.db.ListDomainNames()
		if err != nil {
			return err
		}

		domainNamesListStr := make([]string, 0)
		for _, dn := range domainNamesList {
			domainNamesListStr = append(domainNamesListStr, dn.Name)
		}

		// configure the HTTPs server
		p.server = p.newHttpsServer(domainNamesListStr...)
		p.server.Handler = router
		err = p.server.ListenAndServeTLS("", "")
		return err

	} else {
		// configure the HTTP server
		p.server = p.newHttpServer()
		p.server.Handler = router
		err := p.server.ListenAndServe()
		return err
	}
}

func (p *Proxy) Shutdown() {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	p.server.Shutdown(ctx)
}
