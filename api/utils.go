package api

import (
	"net/http"
	"time"

	"github.com/Dadard29/podman-proxy/models"
)

func (a *Api) GetRuleFromDomainName(domainName string) (models.Rule, error) {
	return a.db.GetRuleFromDomainName(domainName)
}

func (a *Api) GetContainer(containerName string) (models.Container, error) {
	return a.db.GetContainer(containerName)
}

func (a *Api) UpdateDomainNameLive() error {
	return a.db.UpdateDomainNameLive()
}

func (a *Api) ListDomainNames() ([]models.DomainName, error) {
	return a.db.ListDomainNames()
}

func (a *Api) NewNetworkLog(beganAt time.Time, r *http.Request,
	responseLog *models.LogResponseWriter) (*models.NetworkLog, error) {

	duration := time.Since(beganAt)
	responseBody := responseLog.Buf.Bytes()
	netLog, err := models.NewNetworkLogFromRequest(r,
		responseLog.StatusCode, duration, responseBody)

	if err != nil {
		return nil, err
	}

	err = a.db.InsertNetworkLog(netLog)
	if err != nil {
		return nil, err
	}

	return &netLog, nil
}
