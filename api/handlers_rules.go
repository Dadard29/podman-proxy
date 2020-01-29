package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// controller for the /rules endpoint
func rulesHandler(w http.ResponseWriter, r *http.Request) {
	var res response
	var httpStatusCode int

	// set the parameters in the request body
	// set the content-type with the following value: application/json
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// it always returns json data
	w.Header().Add("Content-Type", "application/json")

	var req request
	err = json.Unmarshal(data, &req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		res = response{
			Status:  false,
			Message: "malformed params",
			Rule:    RuleModel{},
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	if r.Method == http.MethodPost {
		res, httpStatusCode = rulesHandlerPost(req)
	} else if r.Method == http.MethodGet {
		res, httpStatusCode = rulesHandlerGet(req)
	} else if r.Method == http.MethodDelete {
		res, httpStatusCode = rulesHandlerDelete(req)
	} else if r.Method == http.MethodPut {
		res, httpStatusCode = rulesHandlerPut(req)
	} else {
		httpStatusCode = http.StatusMethodNotAllowed
		res = response{
			Status:  false,
			Message: "method not allowed",
			Rule:    RuleModel{},
		}
	}

	w.WriteHeader(httpStatusCode)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalln(err)
	}
}

// swagger:route POST /rules rules createRule
//
// Create a rule mapping an hostname and a podman container
//
// This will create a new rule.
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
func rulesHandlerPost(r request) (response, int) {
	// need name, host and port fulfilled

	defaultResponse := response{
		Status:  false,
		Message: "",
		Rule: RuleModel{
			ContainerName: r.ContainerName,
			ContainerHost: r.ContainerHost,
			ContainerPort: r.ContainerPort,
		},
	}

	rule, err := globalApi.CreateRule(r.ContainerName, r.ContainerPort, r.ContainerHost)
	if err != nil {
		log.Println(err)
		defaultResponse.Message = err.Error()
		return defaultResponse, http.StatusInternalServerError
	}

	return response{
		Status:  true,
		Message: "rule created",
		Rule:    rule,
	}, http.StatusOK
}

// swagger:route GET /rules rules getRule
//
// Retrieve a rule details
//
// This will extract rule details from the database
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
func rulesHandlerGet(r request) (response, int) {
	// need host fulfilled

	defaultResponse := response{
		Status:  false,
		Message: "",
		Rule: RuleModel{
			ContainerHost: r.ContainerHost,
		},
	}

	rule, err := globalApi.GetRule(r.ContainerHost)
	if err != nil {
		log.Println(err)
		defaultResponse.Message = err.Error()
		return defaultResponse, http.StatusInternalServerError
	}

	return response{
		Status:  true,
		Message: "rule retrieved",
		Rule:    rule,
	}, http.StatusOK
}

// swagger:route DELETE /rules rules deleteRule
//
// Delete a rule
//
// This will remove a rule
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
func rulesHandlerDelete(r request) (response, int) {
	// need host fulfilled

	defaultResponse := response{
		Status:  false,
		Message: "",
		Rule: RuleModel{
			ContainerHost: r.ContainerHost,
		},
	}

	rule, err := globalApi.DeleteRule(r.ContainerHost)
	if err != nil {
		log.Println(err)
		defaultResponse.Message = err.Error()
		return defaultResponse, http.StatusInternalServerError
	}

	return response{
		Status:  true,
		Message: "rule deleted",
		Rule:    rule,
	}, http.StatusOK
}

// swagger:route PUT /rules rules updateRule
//
// Update a rule concerning a specific container
//
// This will update the IP, the container name, the container port
// and the host, in the rule.
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
func rulesHandlerPut(r request) (response, int) {
	// need host, name and port fulfilled

	defaultResponse := response{
		Status:  false,
		Message: "",
		Rule: RuleModel{
			ContainerName: r.ContainerName,
			ContainerHost: r.ContainerHost,
			ContainerPort: r.ContainerPort,
		},
	}

	rule, err := globalApi.UpdateRule(r.ContainerHost, r.ContainerName, r.ContainerPort)
	if err != nil {
		log.Println(err)
		defaultResponse.Message = err.Error()
		return defaultResponse, http.StatusInternalServerError
	}

	return response{
		Status:  true,
		Message: "rule updated",
		Rule:    rule,
	}, http.StatusOK
}
