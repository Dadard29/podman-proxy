package models

import (
	"io/ioutil"
	"net/http"
	"time"
)

type NetworkLog struct {
	Timestamp time.Time `json:"timestamp"`

	RequestHost   string `json:"request_host"`
	RequestMethod string `json:"request_method"`
	RequestPath   string `json:"request_path"`
	RequestBody   []byte `json:"request_body"`
	RequestArgs   string `json:"request_args"`

	ResponseStatusCode int    `json:"response_status_code"`
	ResponseDuration   int    `json:"response_duration"`
	ResponseBody       []byte `json:"response_body"`
}

func NewNetworkLogFromRequest(req *http.Request,
	responseCode int, responseDuration time.Duration, responseBody []byte) (NetworkLog, error) {

	var body []byte
	if req.Body != nil {
		var err error
		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return NetworkLog{}, err
		}

	} else {
		body = []byte{}
	}

	netLog := NetworkLog{
		Timestamp:     time.Now(),
		RequestHost:   req.Host,
		RequestMethod: req.Method,
		RequestPath:   req.URL.Path,
		RequestBody:   body,
		RequestArgs:   req.URL.RawQuery,

		ResponseStatusCode: responseCode,
		ResponseDuration:   int(responseDuration.Milliseconds()),
		ResponseBody:       responseBody,
	}

	return netLog, nil

}
