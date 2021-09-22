package db

import (
	"context"
	"fmt"

	"github.com/Dadard29/podman-proxy/models"
)

const containerTableName = "containers"

func (db *Db) ListContainers() ([]models.Container, error) {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	rows, err := db.conn.QueryContext(
		ctx,
		fmt.Sprintf("select * from %s", containerTableName),
	)
	if err != nil {
		return nil, err
	}

	out := make([]models.Container, 0)
	for rows.Next() {
		container, err := models.NewContainerFromRow(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, container)
	}

	return out, nil
}

func (db *Db) GetContainer(containerName string) (models.Container, error) {
	var out models.Container

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	row, err := db.conn.QueryContext(
		ctx,
		fmt.Sprintf("select * from %s where name == ?", containerTableName),
		containerName,
	)
	if err != nil {
		return out, err
	}

	for row.Next() {
		return models.NewContainerFromRow(row.Scan)
	}

	return out, fmt.Errorf("container with name %s not found", containerName)
}

func (db *Db) InsertContainer(c *models.Container) error {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	_, err := db.conn.ExecContext(
		ctx,
		fmt.Sprintf(`
		INSERT INTO
		%s(id, name, is_infra, is_in_pod, pod_id, ip_address, status, exposed_port)
		values(?, ?, ?, ?, ?, ?, ?, ?)`, containerTableName),
		c.Id, c.Name, c.IsInfra, c.IsInPod, c.PodId, c.IpAddress, c.Status.String(), 0,
	)
	return err
}

func (db *Db) UpdateContainer(containerName string, ipAddress string, status models.ContainerStatus) error {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	result, err := db.conn.ExecContext(
		ctx,
		fmt.Sprintf("update %s set ip_address = ?, status = ? where name = ?", containerTableName),
		ipAddress, status.String(), containerName,
	)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("update failed: container with name %s not found", containerName)
	}
	return nil
}

func (db *Db) DeleteContainer(containerName string) error {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	result, err := db.conn.ExecContext(
		ctx,
		fmt.Sprintf("delete from %s where name == ?", containerTableName),
		containerName,
	)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("deletion failed: container with name %s not found", containerName)
	}
	return nil
}
