package proxy

import (
	"net/http"

	"github.com/Dadard29/podman-proxy/models"
	"github.com/gorilla/mux"
)

func (p *Proxy) domainNamesHandler(w http.ResponseWriter, r *http.Request) {
	domainNames, err := p.db.ListDomainNames()
	if err != nil {
		p.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		p.WriteErrorJson(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	p.WriteJson(w, &domainNames)
}

func (p *Proxy) domainNameGet(w http.ResponseWriter, r *http.Request, dn string) {
	domainName, err := p.db.GetDomainName(dn)
	if err != nil {
		p.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		p.WriteErrorJson(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	p.WriteJson(w, &domainName)
}

func (p *Proxy) domainNamePost(w http.ResponseWriter, r *http.Request, dn string) {
	err := p.db.InsertDomainName(models.DomainName{
		Name: dn,
	})
	if err != nil {
		p.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		p.WriteErrorJson(w, err)
		return
	}

	domainName, err := p.db.GetDomainName(dn)
	if err != nil {
		p.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		p.WriteErrorJson(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	p.WriteJson(w, &domainName)
}

func (p *Proxy) domainNameDelete(w http.ResponseWriter, r *http.Request, dn string) {
	domainName, err := p.db.GetDomainName(dn)
	if err != nil {
		p.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		p.WriteErrorJson(w, err)
		return
	}

	err = p.db.DeleteDomainName(dn)
	if err != nil {
		p.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		p.WriteErrorJson(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	p.WriteJson(w, &domainName)
}

func (p *Proxy) domainNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dn := vars["dn"]

	if r.Method == http.MethodGet {
		p.domainNameGet(w, r, dn)

	} else if r.Method == http.MethodPost {
		p.domainNamePost(w, r, dn)

	} else if r.Method == http.MethodDelete {
		p.domainNameDelete(w, r, dn)
	}
}
