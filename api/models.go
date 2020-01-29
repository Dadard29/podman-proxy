package api

import "net/http"

type Route struct {
	HttpMethods []string
	Handler func(w http.ResponseWriter, r *http.Request)
}

// swagger:response basicResponse
type response struct {
	// the status of the operation
	Status  bool `json:"Status"`

	// a detailed message about the operation status
	Message string `json:"Message"`

	// the rule model if applicable
	Rule    RuleModel `json:"Rule"`
}

// the request model
type request struct {
	// the hostname to set for the container, always required
	ContainerHost string `json:"ContainerHost"`

	// the name of the container
	ContainerName string `json:"ContainerName"`

	// the port exposed by the container
	ContainerPort int `json:"ContainerPort"`
}
