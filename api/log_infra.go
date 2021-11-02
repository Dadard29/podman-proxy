package api

import (
	"time"

	"github.com/Dadard29/podman-proxy/models"
)

func (a *Api) NewInfraLog() error {
	cpu, err := a.infra.GetCpuUsage(3 * time.Second)
	if err != nil {
		return err
	}
	mem, err := a.infra.GetMemUsage()
	if err != nil {
		return err
	}
	disk, err := a.infra.GetDiskUsage()
	if err != nil {
		return err
	}

	infraLog := models.InfraLog{
		Timestamp: time.Now(),
		Cpu:       cpu,
		Memory:    mem,
		Disk:      disk,
	}

	return a.db.InsertInfraLog(infraLog)
}
