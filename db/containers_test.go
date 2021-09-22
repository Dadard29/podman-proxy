package db_test

import (
	"log"
	"testing"

	"github.com/Dadard29/podman-proxy/models"
)

func TestContainers(t *testing.T) {
	dbService, err := NewTestDb()
	if err != nil {
		t.Error(err)
	}
	defer CleanTestDb()

	// list containers with 0 result
	containersList, err := dbService.ListContainers()
	if err != nil {
		t.Error(err)
	}
	if len(containersList) != 0 {
		t.Errorf("unexpected containers list length: %d", len(containersList))
	}

	// create container
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
	err = dbService.InsertContainer(&container)
	if err != nil {
		t.Error(err)
	}

	// retrieve container
	foundContainer, err := dbService.GetContainer(container.Name)
	if err != nil {
		t.Error(err)
	}

	if foundContainer.Name != container.Name {
		t.Errorf("mismatch: %s != %s", foundContainer.String(), container.String())
	}

	// list containers with 1 result
	containersList, err = dbService.ListContainers()
	if err != nil {
		t.Error(err)
	}
	if len(containersList) != 1 {
		t.Errorf("unexpected containers list length: %d", len(containersList))
	}

	// delete container
	err = dbService.DeleteContainer(container.Name)
	if err != nil {
		t.Error(err)
	}

	// list containers with 0 result
	containersList, err = dbService.ListContainers()
	if err != nil {
		t.Error(err)
	}
	if len(containersList) != 0 {
		t.Errorf("unexpected containers list length: %d", len(containersList))
	}
}

func TestContainersErrors(t *testing.T) {
	dbService, err := NewTestDb()
	if err != nil {
		t.Error(err)
	}
	defer CleanTestDb()

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

	// retrieve container - ERR
	_, err = dbService.GetContainer(container.Name)
	if err == nil {
		t.Error("expected error on get")
	} else {
		log.Println(err)
	}

	// delete container - ERR
	err = dbService.DeleteContainer(container.Name)
	if err == nil {
		t.Error("expected error on delete")
	} else {
		log.Println(err)
	}

	// create container
	err = dbService.InsertContainer(&container)
	if err != nil {
		t.Error(err)
	}

	// (re)create container - ERR
	err = dbService.InsertContainer(&container)
	if err == nil {
		t.Error("expected error on re-creation")
	} else {
		log.Println(err)
	}
}
