package proxy

import "net/http"

func (p *Proxy) ruleListGetHandler(w http.ResponseWriter, r *http.Request) {
	rule, err := p.api.RuleListGet(w, r)

	if err != nil {
		p.WriteErrorJson(w, http.StatusUnauthorized, err)
		return
	}

	p.WriteJson(w, rule)
}

func (p *Proxy) ruleListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.ruleListGetHandler(w, r)
	}
}
