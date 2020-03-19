package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/Dadard29/podman-proxy/api"
	"golang.org/x/crypto/acme/autocert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var globalProxy *Proxy

type Proxy struct {
	exposedApi *api.Api
	httpServer *http.Server
	config     Config
}

const authRoute = "/auth"

// if the given secret is correct, give the associated token
func authHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var payload map[string]string
	err = json.Unmarshal(data, &payload)
	if err != nil {
		log.Println(err)
	}

	secret := payload["Secret"]

	givenToken := globalProxy.config.generateToken(secret)
	if givenToken != globalProxy.config.ProxyToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(givenToken))
	if err != nil {
		log.Println(err)
	}
}

// check if the given token is correct
func checkAuthToken(r *http.Request) bool {
	// the auth key must be in the `Authorization` header, with the value `Bearer <key>`
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return false
	}

	values := strings.Split(authorizationHeader, "Bearer ")
	if len(values) != 2 {
		return false
	}

	authorizationKey := values[1]

	return authorizationKey == globalProxy.config.ProxyToken
}

// send the http request to the required handler
func apiHandler(w http.ResponseWriter, r *http.Request) {
	// check if redirect to auth handler or not
	if r.URL.String() == authRoute {
		authHandler(w, r)
		return
	}

	if !checkAuthToken(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// redirect the request to the correct handler
	for u, h := range globalProxy.exposedApi.GetRoutes() {
		if u == r.URL.String() && checkMethod(h.HttpMethods, r.Method) {
			h.Handler(w, r)
			return
		}
	}

	// no valid handler found for this route and this http method
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte(fmt.Sprintf("%d %s\n", http.StatusNotFound, "custom page not found")))
	if err != nil {
		log.Fatalln(err)
	}
}

// redirect the HTTP traffic to the podman container
func redirectionHandler(w http.ResponseWriter, r *http.Request) {
	splitted := strings.Split(r.Host, ":")
	var requestedHost = splitted[0]

	// retrieve the container ip and port from the request ContainerHost
	rule, err := globalProxy.exposedApi.GetRule(requestedHost)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	containerUrlStr := fmt.Sprintf("http://%s:%d", rule.ContainerIp, rule.ContainerPort)
	containerUrl, err := url.Parse(containerUrlStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	proxyService := httputil.NewSingleHostReverseProxy(containerUrl)

	r.URL.Host = containerUrl.Host
	r.Host = containerUrl.Host
	r.URL.Scheme = containerUrl.Host
	r.Header.Set("X-Forwarded-Host", r.Header.Get("ContainerHost"))

	proxyService.ServeHTTP(w, r)
}

// check if the request is addressed to the API or the container
func mainProxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Host == globalProxy.config.getAddr() {
		apiHandler(w, r)

	} else {
		redirectionHandler(w, r)
	}
}

// proxy constructor
func NewProxy(conf Config) *Proxy {
	//proxyUrl := fmt.Sprintf("%s:%d", conf.ProxyHost, conf.ProxyPort)

	// setup the main proxy route
	http.HandleFunc("/", mainProxyHandler)

	// tls
	manager := autocert.Manager{
		Prompt:          autocert.AcceptTOS,
		Cache:           autocert.DirCache("/srv/https/certificates"),
		HostPolicy:      autocert.HostWhitelist(
			"dadard.fr", "www.dadard.fr",
			"proxy.dadard.fr", "www.proxy.dadard.fr",
			"core.dadard.fr", "www.core.dadard.fr"),
	}

	server := &http.Server{
		Addr:              ":https",
		Handler:           nil,
		TLSConfig:         manager.TLSConfig(),
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
		config:     conf,
	}

	globalProxy = p
	return p
}

// the http server of the proxy
func (p *Proxy) Start() {
	log.Printf("listening on %s...\n", p.httpServer.Addr)
	err := p.httpServer.ListenAndServeTLS("", "")

	if err != nil {
		log.Fatalln(err)
	}
}
