package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type response struct {
	Status  bool
	Message string
	Rule    RuleModel
}

type request struct {
	ContainerHost string
	ContainerName string
	ContainerPort int
}

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
		w.WriteHeader(http.StatusMethodNotAllowed)
		res = response{
			Status:  false,
			Message: "method not allowed",
			Rule:    RuleModel{},
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Fatalln(err)
		}
		return

	}

	w.WriteHeader(httpStatusCode)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalln(err)
	}
}

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
