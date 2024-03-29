package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Dadard29/podman-proxy/api"
	"github.com/Dadard29/podman-proxy/web"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
)

type Proxy struct {
	Upgrader *Upgrader

	config config
	logger *log.Logger
	server *http.Server
	router *mux.Router
	api    *api.Api

	Ctx    *context.Context
	Cancel *context.CancelFunc
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
	domainNames = append(domainNames, p.config.proxyHost)
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
		ErrorLog:  log.New(os.Stdout, "proxy-error ", log.Ldate|log.Ltime),
	}

	return server
}

func (p *Proxy) newHttpServer() *http.Server {

	// server
	addr := p.getAddrFromProxyPort()
	return &http.Server{
		Addr:     addr,
		ErrorLog: log.New(os.Stdout, "proxy-error ", log.Ldate|log.Ltime),
	}
}

func NewProxy() (*Proxy, error) {
	config, err := newConfigFromEnv()
	if err != nil {
		return nil, err
	}

	logger := log.New(log.Default().Writer(), "proxy ", log.Default().Flags())
	upgrader := NewUpgrader(config.upgraderPort)

	proxyApi, err := api.NewApi(config.dbPath)
	if err != nil {
		return nil, err
	}

	proxy := &Proxy{
		config:   config,
		logger:   logger,
		server:   nil,
		router:   nil,
		Upgrader: upgrader,
		api:      proxyApi,

		Ctx:    nil,
		Cancel: nil,
	}

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/auth", proxy.authHandler).Methods(http.MethodPost, http.MethodGet)
	apiRouter.HandleFunc("/rule", proxy.ruleListHandler).Methods(http.MethodGet)
	apiRouter.HandleFunc("/rule/{dn}", proxy.ruleHandler).Methods(http.MethodGet, http.MethodPost, http.MethodDelete)
	apiRouter.HandleFunc("/domain-name", proxy.domainNameListHandler).Methods(http.MethodGet, http.MethodPut)
	apiRouter.HandleFunc("/domain-name/{dn}", proxy.domainNameHandler).Methods(http.MethodGet, http.MethodPost, http.MethodDelete)
	apiRouter.HandleFunc("/container", proxy.containerListHandler).Methods(http.MethodGet, http.MethodPut)
	apiRouter.HandleFunc("/container/{container}", proxy.containerHandler).Methods(http.MethodGet, http.MethodPost)
	apiRouter.HandleFunc("/container-sync", proxy.containerSyncHandler).Methods(http.MethodGet, http.MethodPost)

	apiRouter.Use(proxy.authMiddleware)
	apiRouter.Use(proxy.dbLoggingMiddleware)

	webRouter := router.PathPrefix("/").Subrouter()
	web.RegisterRoutes(webRouter)

	proxy.router = router

	if proxy.config.debug {
		proxy.logger.Println("WARNING: debug on")
		proxy.logger.Println("WARNING: change the debug parameter for production use")
	}

	return proxy, nil
}

func (p *Proxy) Serve(withTLS bool) error {

	router := mux.NewRouter()
	router.PathPrefix("/").HandlerFunc(p.switcher)
	// router.Use(p.dbLoggingMiddleware)

	ctx, cancel := context.WithCancel(context.Background())
	p.Cancel = &cancel
	p.Ctx = &ctx

	p.api.UpdateDomainNameLive()

	p.logger.Printf("Starting proxy server on %s...\n", p.Host())

	if withTLS {
		// get the list of domain names to register
		domainNamesList, err := p.api.ListDomainNames()
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
		if err != nil && err != http.ErrServerClosed {
			return err
		}

	} else {
		// configure the HTTP server
		p.server = p.newHttpServer()
		p.server.Handler = router
		err := p.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

func (p *Proxy) Shutdown() {

	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		p.Ctx = nil
		p.Cancel = nil

		cancel()
	}()

	p.logger.Println("shutting down server...")
	err := p.server.Shutdown(ctx)
	if err != nil {
		p.logger.Println("WARNING:", err)
	}
}

func (p *Proxy) NewInfraLog() error {
	return p.api.NewInfraLog()
}

func (p *Proxy) Close() {
	p.logger.Println("closing connections...")
	p.api.Close()
}
