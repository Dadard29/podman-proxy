package proxy

import (
	"net/http"

	"github.com/Dadard29/podman-proxy/models"
	"github.com/gorilla/mux"
)

func (p *Proxy) domainNameListUpdate(w http.ResponseWriter, r *http.Request) {

	p.WriteMessageJson(w, "restarting proxy...")

	(*p.Cancel)()
}

func (p *Proxy) domainNameListGet(w http.ResponseWriter, r *http.Request) {
	domainNames, err := p.db.ListDomainNames()
	if err != nil {
		p.logger.Println(err)
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	p.WriteJson(w, &domainNames)
}

func (p *Proxy) domainNamesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.domainNameListGet(w, r)

	} else if r.Method == http.MethodPut {
		p.domainNameListUpdate(w, r)

	}
}

func (p *Proxy) domainNameGet(w http.ResponseWriter, r *http.Request, dn string) {
	domainName, err := p.db.GetDomainName(dn)
	if err != nil {
		p.WriteErrorJson(w, http.StatusNotFound, err)
		return
	}

	p.WriteJson(w, &domainName)
}

func (p *Proxy) domainNamePost(w http.ResponseWriter, r *http.Request, dn string) {
	err := p.db.InsertDomainName(models.DomainName{
		Name: dn,
	})
	if err != nil {
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	domainName, err := p.db.GetDomainName(dn)
	if err != nil {
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	p.WriteJson(w, &domainName)
}

func (p *Proxy) domainNameDelete(w http.ResponseWriter, r *http.Request, dn string) {
	domainName, err := p.db.GetDomainName(dn)
	if err != nil {
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	err = p.db.DeleteDomainName(dn)
	if err != nil {
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	p.WriteJson(w, &domainName)
}

func (p *Proxy) domainNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dn := vars["dn"]

	if r.Method == http.MethodGet {
		p.domainNameGet(w, r, dn)

	} else if r.Method == http.MethodPost {
		p.domainNamePost(w, r, dn)

	} else if r.Method == http.MethodPut {

	} else if r.Method == http.MethodDelete {
		p.domainNameDelete(w, r, dn)
	}
}
