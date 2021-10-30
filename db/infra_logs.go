package db

import (
	"context"
	"fmt"

	"github.com/Dadard29/podman-proxy/models"
)

const infraLogTableName = "infra_logs"

func (db *Db) InsertInfraLog(infraLog models.InfraLog) error {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	_, err := db.conn.ExecContext(
		ctx,
		fmt.Sprintf(`
		INSERT INTO %s(timestamp, cpu, memory, disk)
		values(?, ?, ?, ?)
		`, infraLogTableName),
		infraLog.Timestamp.Unix(),
		infraLog.Cpu,
		infraLog.Memory,
		infraLog.Disk,
	)
	return err
}
