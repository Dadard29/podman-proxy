package api

import (
	"fmt"
	"net/http"

	"github.com/Dadard29/podman-proxy/models"
)

// Retrieve all existing rules
func (a *Api) RuleListGet(w http.ResponseWriter, r *http.Request) (*[]models.Rule, error) {
	rules, err := a.db.ListRules()
	if err != nil {
		return nil, err
	}

	return &rules, nil
}

// Retrieve an existing rule
func (a *Api) RuleGet(w http.ResponseWriter, r *http.Request, dn string) (*models.Rule, error) {

	rule, err := a.db.GetRuleFromDomainName(dn)
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

// Create a new rule
func (a *Api) RulePost(w http.ResponseWriter, r *http.Request, dn string) (*models.Rule, error) {
	containerName := r.URL.Query().Get("containerName")

	// check if container has a valid exposedPort
	container, err := a.db.GetContainer(containerName)
	if err != nil {
		return nil, err
	}

	if container.ExposedPort == 0 {
		err := fmt.Errorf("container %s has no valid exposedPort set", container.Name)
		return nil, err
	}

	err = a.db.InsertRule(dn, containerName)
	if err != nil {
		return nil, err
	}

	rule, err := a.db.GetRuleFromDomainName(dn)
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

// Delete a rule
func (a *Api) RuleDelete(w http.ResponseWriter, r *http.Request, dn string) (*models.Rule, error) {
	rule, err := a.db.GetRuleFromDomainName(dn)
	if err != nil {
		return nil, err
	}

	err = a.db.DeleteRuleFromDomainName(dn)
	if err != nil {
		return nil, err
	}

	return &rule, nil
}
