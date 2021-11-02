package proxy

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (p *Proxy) domainNameGetHandler(w http.ResponseWriter, r *http.Request, dn string) {
	domainName, err := p.api.DomainNameGet(w, r, dn)

	if err != nil {
		p.WriteErrorJson(w, http.StatusUnauthorized, err)
		return
	}

	p.WriteJson(w, domainName)
}

func (p *Proxy) domainNamePostHandler(w http.ResponseWriter, r *http.Request, dn string) {
	domainName, err := p.api.DomainNamePost(w, r, dn)

	if err != nil {
		p.WriteErrorJson(w, http.StatusUnauthorized, err)
		return
	}

	p.WriteJson(w, domainName)
}

func (p *Proxy) domainNameDeleteHandler(w http.ResponseWriter, r *http.Request, dn string) {
	domainName, err := p.api.DomainNameDelete(w, r, dn)

	if err != nil {
		p.WriteErrorJson(w, http.StatusUnauthorized, err)
		return
	}

	p.WriteJson(w, domainName)
}

func (p *Proxy) domainNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dn := vars["dn"]

	if r.Method == http.MethodGet {
		p.domainNameGetHandler(w, r, dn)

	} else if r.Method == http.MethodPost {
		p.domainNamePostHandler(w, r, dn)

	} else if r.Method == http.MethodDelete {
		p.domainNameDeleteHandler(w, r, dn)
	}
}
