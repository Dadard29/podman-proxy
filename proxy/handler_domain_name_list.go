package proxy

import "net/http"

func (p *Proxy) domainNameListGetHandler(w http.ResponseWriter, r *http.Request) {
	domainNameList, err := p.api.DomainNameListGet(w, r)
	if err != nil {
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	p.WriteJson(w, domainNameList)
}

func (p *Proxy) domainNameListUpdateHandler(w http.ResponseWriter, r *http.Request) {

	p.WriteMessageJson(w, "restarting proxy...")
	(*p.Cancel)()
}

func (p *Proxy) domainNameListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.domainNameListGetHandler(w, r)

	} else if r.Method == http.MethodPut {
		p.domainNameListUpdateHandler(w, r)

	}
}
