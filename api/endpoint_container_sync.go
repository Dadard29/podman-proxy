package api

import (
	"net/http"
	"time"

	"github.com/Dadard29/podman-proxy/models"
)

type containerSyncTime struct {
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

func (a *Api) ContainerSyncGet(w http.ResponseWriter, r *http.Request) (*containerSyncTime, error) {
	lastUpdatedAt, err := a.db.GetContainerLastUpdatedAt()
	if err != nil {
		return nil, err
	}

	return &containerSyncTime{
		LastUpdatedAt: *lastUpdatedAt,
	}, nil
}

// Update the database with all existing containersPodman
func (a *Api) ContainerSyncPost(w http.ResponseWriter, r *http.Request) (*[]models.Container, error) {
	containersDb, err := a.db.ListContainers()
	if err != nil {
		return nil, err
	}

	containersPodman, err := a.podman.ListContainers()
	if err != nil {
		return nil, err
	}

	for _, containerPodman := range containersPodman {
		found := false

		// updating the database with the new IPs of the containers
		for _, containerDb := range containersDb {
			if containerPodman.Name == containerDb.Name {
				err := a.db.UpdateContainerIpStatus(containerPodman.Name, containerPodman.IpAddress, containerPodman.Status)
				if err != nil {
					a.logger.Println(err)
				}
				found = true
				break
			}
		}

		// creating in database the newly created containers
		if !found {
			err := a.db.InsertContainer(containerPodman)
			if err != nil {
				a.logger.Println(err)
			}
		}
	}

	// deleting the non-existant containers in database
	for _, containerDb := range containersDb {
		found := false

		for _, containerPodman := range containersPodman {
			if containerPodman.Name == containerDb.Name {
				found = true
			}
		}

		if !found {
			err := a.db.DeleteContainer(containerDb.Name)
			if err != nil {
				a.logger.Println(err)
			}
		}
	}

	// retrieve the new database containers
	containersDb, err = a.db.ListContainers()
	if err != nil {
		return nil, err
	}

	return &containersDb, nil
}
