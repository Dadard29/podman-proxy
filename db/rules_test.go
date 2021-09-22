package db_test

import (
	"log"
	"testing"

	"github.com/Dadard29/podman-proxy/models"
)

func TestRules(t *testing.T) {
	dbService, err := NewTestDb()
	if err != nil {
		t.Error(err)
	}
	defer CleanTestDb()

	// init foreign refs
	dn := models.DomainName{
		Name: "host.com",
	}
	dbService.InsertDomainName(dn)

	container := models.Container{
		Id:          "id",
		Name:        "container",
		IsInfra:     false,
		IsInPod:     false,
		PodId:       "",
		IpAddress:   "10.10.10.10",
		ExposedPort: 0,
		Status:      models.NewContainerStatus("running"),
	}
	dbService.InsertContainer(&container)

	// list rules with 0 results
	ruleList, err := dbService.ListRules()
	if err != nil {
		t.Error(err)
	}
	if len(ruleList) != 0 {
		t.Errorf("unexpected rule list length: %d", len(ruleList))
	}

	// create rule
	err = dbService.InsertRule(dn.Name, container.Name)
	if err != nil {
		t.Error(err)
	}

	// get rule
	foundRule, err := dbService.GetRuleFromDomainName(dn.Name)
	if err != nil {
		t.Error(err)
	}
	if foundRule.DomainName != dn.Name || foundRule.ContainerName != container.Name {
		t.Errorf("mismatch: %s != (%s -> %s)", foundRule.String(), dn.Name, container.Name)
	}

	// list rules with 1 result
	ruleList, err = dbService.ListRules()
	if err != nil {
		t.Error(err)
	}
	if len(ruleList) != 1 {
		t.Errorf("unexpected rule list length: %d", len(ruleList))
	}

	// delete rule
	err = dbService.DeleteRuleFromDomainName(dn.Name)
	if err != nil {
		t.Error(err)
	}

	// list rules with 0 result
	ruleList, err = dbService.ListRules()
	if err != nil {
		t.Error(err)
	}
	if len(ruleList) != 0 {
		t.Errorf("unexpected rule list length: %d", len(ruleList))
	}
}

func TestRulesErrors(t *testing.T) {
	dbService, err := NewTestDb()
	if err != nil {
		t.Error(err)
	}
	defer CleanTestDb()

	// init foreign refs
	dn := models.DomainName{
		Name: "host.com",
	}
	container := models.Container{
		Id:          "id",
		Name:        "container",
		IsInfra:     false,
		IsInPod:     false,
		PodId:       "",
		IpAddress:   "10.10.10.10",
		ExposedPort: 0,
		Status:      models.NewContainerStatus("running"),
	}

	// get rule - ERR
	_, err = dbService.GetRuleFromDomainName(dn.Name)
	if err == nil {
		t.Error("expected error on retrieve")
	} else {
		log.Println(err)
	}

	// delete rule - ERR
	err = dbService.DeleteRuleFromDomainName(dn.Name)
	if err == nil {
		t.Error("expected error on deletion")
	} else {
		log.Println(err)
	}

	// create rule without foreign key - ERR
	err = dbService.InsertRule(dn.Name, container.Name)
	if err == nil {
		t.Error("expected error on creation")
	} else {
		log.Println(err)
	}

	// init dn
	dbService.InsertDomainName(dn)

	// create rule without foreign key - ERR
	err = dbService.InsertRule(dn.Name, container.Name)
	if err == nil {
		t.Error("expected error on creation")
	} else {
		log.Println(err)
	}

	// init container
	dbService.InsertContainer(&container)

	// create rule
	err = dbService.InsertRule(dn.Name, container.Name)
	if err != nil {
		t.Error(err)
	}

	// (re)create rule - ERR
	err = dbService.InsertRule(dn.Name, container.Name)
	if err == nil {
		t.Error("expected error on re-creation")
	} else {
		log.Println(err)
	}
}
