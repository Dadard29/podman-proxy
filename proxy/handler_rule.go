package proxy

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Retrieve all existing rules
func (p *Proxy) rulesHandler(w http.ResponseWriter, r *http.Request) {
	rules, err := p.db.ListRules()
	if err != nil {
		p.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	p.WriteJson(w, &rules)
}

// Retrieve an existing rule
func (p *Proxy) ruleGet(w http.ResponseWriter, r *http.Request, dn string) {

	rule, err := p.db.GetRuleFromDomainName(dn)
	if err != nil {
		p.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	p.WriteJson(w, &rule)
}

// Create a new rule
func (p *Proxy) rulePost(w http.ResponseWriter, r *http.Request, dn string) {
	containerName := r.URL.Query().Get("containerName")
	err := p.db.InsertRule(dn, containerName)
	if err != nil {
		p.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rule, err := p.db.GetRuleFromDomainName(dn)
	if err != nil {
		p.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	p.WriteJson(w, &rule)
}

// Delete a rule
func (p *Proxy) ruleDelete(w http.ResponseWriter, r *http.Request, dn string) {
	rule, err := p.db.GetRuleFromDomainName(dn)
	if err != nil {
		p.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = p.db.DeleteRuleFromDomainName(dn)
	if err != nil {
		p.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	p.WriteJson(w, &rule)
}

// Main rule handler
func (p *Proxy) ruleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dn := vars["dn"]

	if r.Method == http.MethodGet {
		p.ruleGet(w, r, dn)

	} else if r.Method == http.MethodPost {
		p.rulePost(w, r, dn)

	} else if r.Method == http.MethodDelete {
		p.ruleDelete(w, r, dn)
	}
}
