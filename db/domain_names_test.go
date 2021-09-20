package db_test

import (
	"log"
	"testing"

	"github.com/Dadard29/podman-proxy/models"
)

func TestDomainNames(t *testing.T) {
	dbService, err := NewTestDb()
	if err != nil {
		t.Error(err)
	}
	defer CleanTestDb()

	// list dn
	dnList, err := dbService.ListDomainNames()
	if err != nil {
		t.Error(err)
	}
	if len(dnList) != 0 {
		t.Errorf("unexpected dn list: %d", len(dnList))
	}

	// create dn
	dn := models.DomainName{
		Name: "host.com",
	}
	err = dbService.InsertDomainName(dn)
	if err != nil {
		t.Error(err)
	}

	// get dn
	foundDn, err := dbService.GetDomainName(dn.Name)
	if err != nil {
		t.Error(err)
	}
	if foundDn.Name != dn.Name {
		t.Errorf("mismatch: %s != %s", foundDn.Name, dn.Name)
	}

	// list dn
	dnList, err = dbService.ListDomainNames()
	if err != nil {
		t.Error(err)
	}
	if len(dnList) != 1 {
		t.Errorf("unexpected dn list: %d", len(dnList))
	}

	// delete dn
	err = dbService.DeleteDomainName(dn.Name)
	if err != nil {
		t.Error(err)
	}

	// list dn with 0 results
	dnList, err = dbService.ListDomainNames()
	if err != nil {
		t.Error(err)
	}
	if len(dnList) != 0 {
		t.Errorf("unexpected dn list: %d", len(dnList))
	}
}

func TestDomainNamesErrors(t *testing.T) {
	dbService, err := NewTestDb()
	if err != nil {
		t.Error(err)
	}
	defer CleanTestDb()

	dn := models.DomainName{
		Name: "host.com",
	}

	// get dn - ERR
	_, err = dbService.GetDomainName(dn.Name)
	if err == nil {
		t.Error("expected error on get")
	} else {
		log.Println(err)
	}

	// delete dn - ERR
	err = dbService.DeleteDomainName(dn.Name)
	if err == nil {
		t.Error("expected error on delete")
	} else {
		log.Println(err)
	}

	// create dn
	err = dbService.InsertDomainName(dn)
	if err != nil {
		t.Error(err)
	}

	// (re)create dn - ERR
	err = dbService.InsertDomainName(dn)
	if err == nil {
		t.Error("expected error on re-creation")
	} else {
		log.Println(err)
	}
}
