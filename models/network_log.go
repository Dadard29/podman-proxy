package models

import (
	"io/ioutil"
	"net/http"
	"strconv"
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

func NewNetworkLog(scan func(dest ...interface{}) error) (NetworkLog, error) {
	var timestampStr string
	var requestHost string
	var requestMethod string
	var requestPath string
	var requestBody []byte
	var requestArgs string
	var responseStatusCode int
	var responseDuration int
	var responseBody []byte

	err := scan(&timestampStr, &requestHost, &requestMethod, &requestPath, &requestBody,
		&requestArgs, &responseStatusCode, &responseDuration, &responseBody)

	if err != nil {
		return NetworkLog{}, err
	}

	timestampInt, err := strconv.Atoi(timestampStr)
	if err != nil {
		return NetworkLog{}, err
	}
	timestamp := time.Unix(int64(timestampInt), 0)

	return NetworkLog{
		Timestamp:          timestamp,
		RequestHost:        requestHost,
		RequestMethod:      requestMethod,
		RequestPath:        requestPath,
		RequestBody:        requestBody,
		RequestArgs:        requestArgs,
		ResponseStatusCode: responseStatusCode,
		ResponseDuration:   responseDuration,
		ResponseBody:       requestBody,
	}, nil

}
