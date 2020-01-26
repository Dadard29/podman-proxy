package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type responseRulesList struct {
	Status bool
	Message string
	Rule []RuleModel
}

func rulesListHandler(w http.ResponseWriter, r *http.Request) {
	var res responseRulesList
	var httpStatusCode int
	var err error

	// it always returns json data
	w.Header().Add("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		res, httpStatusCode = rulesHandlerGetList()
	} else {
		httpStatusCode = http.StatusMethodNotAllowed
		res = responseRulesList{
			Status:  false,
			Message: "method not allowed",
			Rule:    []RuleModel{},
		}
	}

	w.WriteHeader(httpStatusCode)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalln(err)
	}
}


func rulesHandlerGetList() (responseRulesList, int) {
	rulesList, err := globalApi.ListRules()
	if err != nil {
		log.Print(err)
		return responseRulesList{
			Status:  false,
			Message: err.Error(),
			Rule:    nil,
		}, http.StatusInternalServerError
	}
	return responseRulesList{
		Status:  true,
		Message: "rule list retrieved",
		Rule:    rulesList,
	}, http.StatusOK
}
