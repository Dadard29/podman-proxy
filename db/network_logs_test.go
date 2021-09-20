package db_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/Dadard29/podman-proxy/models"
)

func TestNeworkLogs(t *testing.T) {
	dbService, err := NewTestDb()
	if err != nil {
		t.Error(err)
	}
	defer CleanTestDb()

	req, err := http.NewRequest(http.MethodGet, "http://host.com", nil)
	if err != nil {
		t.Error(err)
	}
	responseBody := []byte{}
	responseDuration := time.Duration(10)
	responseCode := 500
	netLog, err := models.NewNetworkLogFromRequest(req, responseCode, responseDuration, responseBody)
	if err != nil {
		t.Error(err)
	}

	err = dbService.InsertNetworkLog(netLog)
	if err != nil {
		t.Error(err)
	}
}
