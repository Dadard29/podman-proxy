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

// controller for the /rules/list endpoint
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

// swagger:route GET /rules/list rules list listRule
//
// List the stored rules
//
// This will extract all rules details from the database
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: basicResponse
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
