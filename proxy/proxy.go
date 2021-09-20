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

func getConfigFromEnv() (config, error) {
	proxyHost := os.Getenv(envProxyHost)
	dbPath := os.Getenv(envDbPath)
	addr := os.Getenv(envAddr)

	if proxyHost == "" {
		return config{}, fmt.Errorf("env variable not set: %s", envProxyHost)
	}
	if dbPath == "" {
		return config{}, fmt.Errorf("env variable not set: %s", envDbPath)
	}
	if addr == "" {
		return config{}, fmt.Errorf("env variable not set: %s", envAddr)
	}

	proxyConfig := config{
		proxyHost: proxyHost,
		dbPath:    dbPath,
		addr:      addr,
	}

	return proxyConfig, nil
}

func newHttpsServer(addr string, domainNames ...string) *http.Server {
	// tls
	manager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("/srv/https/certificates"),
		HostPolicy: autocert.HostWhitelist(domainNames...),
	}

	server := &http.Server{
		Addr:      addr,
		TLSConfig: manager.TLSConfig(),
		ErrorLog:  log.New(os.Stdout, "podman-proxy-error ", log.Ldate|log.Ltime),
	}

	return server
}

func newHttpServer(addr string) *http.Server {
	return &http.Server{
		Addr:     addr,
		ErrorLog: log.New(os.Stdout, "podman-proxy-error ", log.Ldate|log.Ltime),
	}
}

func NewProxy() (*Proxy, error) {
	config, err := getConfigFromEnv()
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

	p.logger.Printf("Starting proxy server on %s...\n", p.config.addr)

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
		p.server = newHttpsServer(p.config.addr, domainNamesListStr...)
		p.server.Handler = router
		err = p.server.ListenAndServeTLS("", "")
		return err

	} else {
		// configure the HTTP server
		p.server = newHttpServer(p.config.addr)
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
