package db

import (
	"github.com/Dadard29/podman-proxy/api"
	"testing"
)

func TestDbConnector(t *testing.T) {
	a := api.NewApi()
	if ! a.PingDb() {
		t.Errorf("db down")
	}
}

func TestDbGet(t *testing.T) {
	a := api.NewApi()
	_, err := a.GetRule("server-host")
	if err != nil {
		t.Error(err)
	}
}
