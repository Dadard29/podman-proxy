package models

import (
	"strconv"
	"time"
)

type InfraLog struct {
	Timestamp time.Time `json:"timestamp"`

	Cpu    float32 `json:"cpu"`
	Memory float32 `json:"memory"`
	Disk   float32 `json:"disk"`
}

func NewInfraLog(scan func(dest ...interface{}) error) (InfraLog, error) {
	var timestampStr string
	var cpu float32
	var memory float32
	var disk float32

	err := scan(&timestampStr, &cpu, &memory, &disk)
	if err != nil {
		return InfraLog{}, err
	}

	timestampInt, err := strconv.Atoi(timestampStr)
	if err != nil {
		return InfraLog{}, err
	}
	timestamp := time.Unix(int64(timestampInt), 0)

	return InfraLog{
		Timestamp: timestamp,
		Cpu:       cpu,
		Memory:    memory,
		Disk:      disk,
	}, nil
}
