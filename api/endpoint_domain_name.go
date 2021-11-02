package api

import (
	"net/http"

	"github.com/Dadard29/podman-proxy/models"
)

func (a *Api) DomainNameListGet(w http.ResponseWriter, r *http.Request) (*[]models.DomainName, error) {
	domainNames, err := a.db.ListDomainNames()
	if err != nil {
		return nil, err
	}

	return &domainNames, nil
}

func (a *Api) DomainNameGet(w http.ResponseWriter, r *http.Request, dn string) (*models.DomainName, error) {
	domainName, err := a.db.GetDomainName(dn)
	if err != nil {
		return nil, err
	}

	return &domainName, nil
}

func (a *Api) DomainNamePost(w http.ResponseWriter, r *http.Request, dn string) (*models.DomainName, error) {
	err := a.db.InsertDomainName(models.DomainName{
		Name: dn,
	})
	if err != nil {
		return nil, err
	}

	domainName, err := a.db.GetDomainName(dn)
	if err != nil {
		return nil, err
	}

	return &domainName, nil
}

func (a *Api) DomainNameDelete(w http.ResponseWriter, r *http.Request, dn string) (*models.DomainName, error) {
	domainName, err := a.db.GetDomainName(dn)
	if err != nil {
		return nil, err
	}

	err = a.db.DeleteDomainName(dn)
	if err != nil {
		return nil, err
	}

	return &domainName, nil
}
