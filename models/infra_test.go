package models_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Dadard29/podman-proxy/models"
)

func TestInfraCpu(t *testing.T) {
	infra := models.NewInfra()
	cpu, err := infra.GetCpuUsage(5 * time.Second)
	if err != nil {
		t.Errorf("failed getting CPU stats: %v", err)
		return
	}

	fmt.Println(cpu)
}

func TestInfraMem(t *testing.T) {
	infra := models.NewInfra()
	mem, err := infra.GetMemUsage()
	if err != nil {
		t.Errorf("failed getting virtual mem: %v", err)
	}
	fmt.Println(mem)
}

func TestInfraDisk(t *testing.T) {
	infra := models.NewInfra()
	disk, err := infra.GetDiskUsage()
	if err != nil {
		t.Errorf("failed getting disk usage: %v", err)
	}
	fmt.Println(disk)
}
