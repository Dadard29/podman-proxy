package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

func (p *Proxy) WriteJson(w http.ResponseWriter, i interface{}) {
	res, _ := json.MarshalIndent(i, "", "    ")
	w.Write(res)
}

// Redirect the request to the container associated with the domain name used
func (p *Proxy) redirectToContainer(w http.ResponseWriter, r *http.Request) {
	splitted := strings.Split(r.Host, ":")
	dn := splitted[0]

	rule, err := p.db.GetRuleFromDomainName(dn)
	if err != nil {
		p.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	container, err := p.db.GetContainer(rule.ContainerName)
	if err != nil {
		p.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	containerUrlStr := fmt.Sprintf("http://%s:%d", container.IpAddress, container.ExposedPort)
	containerUrl, err := url.Parse(containerUrlStr)
	if err != nil {
		p.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(containerUrl)
	r.URL.Host = containerUrl.Host
	r.Host = container.Name

	reverseProxy.ServeHTTP(w, r)
}

// Uses the built-in proxy routes
func (p *Proxy) redirectToProxyApi(w http.ResponseWriter, r *http.Request) {
	var matcher mux.RouteMatch
	if check := p.router.Match(r, &matcher); check {

		newReq := mux.SetURLVars(r, matcher.Vars)
		matcher.Handler.ServeHTTP(w, newReq)

	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// Main proxy entrypoint
func (p *Proxy) switcher(w http.ResponseWriter, r *http.Request) {
	if r.Host == p.config.proxyHost {
		p.redirectToProxyApi(w, r)

	} else {
		p.redirectToContainer(w, r)

	}
}