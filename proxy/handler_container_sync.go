package proxy

import (
	"net/http"
	"time"
)

func (p *Proxy) containerSyncGet(w http.ResponseWriter, r *http.Request) {
	lastUpdatedAt, err := p.db.GetContainerLastUpdatedAt()
	if err != nil {
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	res := struct {
		LastUpdatedAt time.Time `json:"last_updated_at"`
	}{
		LastUpdatedAt: *lastUpdatedAt,
	}

	p.WriteJson(w, &res)
}

// Update the database with all existing containersPodman
func (p *Proxy) containerSyncPost(w http.ResponseWriter, r *http.Request) {
	containersDb, err := p.db.ListContainers()
	if err != nil {
		p.logger.Println(err)
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	containersPodman, err := p.podman.ListContainers()
	if err != nil {
		p.logger.Println(err)
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	for _, containerPodman := range containersPodman {
		found := false

		// updating the database with the new IPs of the containers
		for _, containerDb := range containersDb {
			if containerPodman.Name == containerDb.Name {
				err := p.db.UpdateContainerIpStatus(containerPodman.Name, containerPodman.IpAddress, containerPodman.Status)
				if err != nil {
					p.logger.Println(err)
				}
				found = true
				break
			}
		}

		// creating in database the newly created containers
		if !found {
			err := p.db.InsertContainer(containerPodman)
			if err != nil {
				p.logger.Println(err)
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
			err := p.db.DeleteContainer(containerDb.Name)
			if err != nil {
				p.logger.Println(err)
			}
		}
	}

	// retrieve the new database containers
	containersDb, err = p.db.ListContainers()
	if err != nil {
		if err != nil {
			p.logger.Println(err)
			p.WriteErrorJson(w, http.StatusInternalServerError, err)
			return
		}
	}

	p.WriteJson(w, &containersDb)

}

func (p *Proxy) containerSyncHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.containerSyncGet(w, r)

	} else if r.Method == http.MethodPost {
		p.containerSyncPost(w, r)

	}
}
