package db

import (
	"github.com/Dadard29/podman-proxy/api"
	"testing"
)

var a = api.NewApi()

func TestDbConnector(t *testing.T) {
	if !a.PingDb() {
		t.Errorf("db down")
	}
}

func TestDbCreateKO(t *testing.T) {
	// podman container does not exists
	_, err := a.CreateRule("server", 8000, "server-host")
	if err == nil {
		t.Error("no error raised")
	}
}

func TestDbGetKO(t *testing.T) {
	// rule does not exist
	_, err := a.GetRule("server-host")
	if err == nil {
		t.Error("no error raised")
	}
}
