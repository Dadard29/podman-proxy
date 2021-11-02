package api

import (
	"log"

	"github.com/Dadard29/podman-proxy/db"
	"github.com/Dadard29/podman-proxy/models"
	"github.com/Dadard29/podman-proxy/podman"
)

type Api struct {
	logger *log.Logger
	db     *db.Db
	podman *podman.PodmanRuntime
	infra  models.Infra
}

func (a *Api) Close() {
	a.db.Close()
	a.podman.Stop()
}

func NewApi(dbPath string) (*Api, error) {
	logger := log.New(log.Default().Writer(), "api ", log.Default().Flags())

	apiDb, err := db.NewDb(dbPath)
	if err != nil {
		return nil, err
	}

	infra := models.NewInfra()

	runtime, err := podman.NewPodmanRuntime()
	if err != nil {
		return nil, err
	}

	return &Api{
		logger: logger,
		db:     apiDb,
		podman: runtime,
		infra:  infra,
	}, nil
}
