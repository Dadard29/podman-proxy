package proxy

import (
	"math"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type Infra struct {
}

func NewInfra() Infra {
	return Infra{}
}

func (i Infra) GetCpuUsage(d time.Duration) (float32, error) {
	stat, err := cpu.Percent(d, false)
	if err != nil {
		return 0, err
	}

	return float32(math.Round(stat[0])), nil
}

func (i Infra) GetMemUsage() (float32, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}

	return float32(v.UsedPercent), nil

}

func (i Infra) GetDiskUsage() (float32, error) {
	usage, err := disk.Usage("/")
	if err != nil {
		return 0, err
	}

	return float32(usage.UsedPercent), nil
}
