package proxy

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (p *Proxy) ruleGetHandler(w http.ResponseWriter, r *http.Request, dn string) {
	rule, err := p.api.RuleGet(w, r, dn)

	if err != nil {
		p.WriteErrorJson(w, http.StatusUnauthorized, err)
		return
	}

	p.WriteJson(w, rule)
}

func (p *Proxy) rulePostHandler(w http.ResponseWriter, r *http.Request, dn string) {
	rule, err := p.api.RulePost(w, r, dn)

	if err != nil {
		p.WriteErrorJson(w, http.StatusUnauthorized, err)
		return
	}

	p.WriteJson(w, rule)
}

func (p *Proxy) ruleDeleteHandler(w http.ResponseWriter, r *http.Request, dn string) {
	rule, err := p.api.RuleDelete(w, r, dn)

	if err != nil {
		p.WriteErrorJson(w, http.StatusUnauthorized, err)
		return
	}

	p.WriteJson(w, rule)
}

// Main rule handler
func (p *Proxy) ruleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dn := vars["dn"]

	if r.Method == http.MethodGet {
		p.ruleGetHandler(w, r, dn)

	} else if r.Method == http.MethodPost {
		p.rulePostHandler(w, r, dn)

	} else if r.Method == http.MethodDelete {
		p.ruleDeleteHandler(w, r, dn)
	}
}
