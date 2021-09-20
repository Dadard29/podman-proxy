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
		container, err := models.NewContainer(rows)
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
		return models.NewContainer(row)
	}

	return out, fmt.Errorf("container with name %s not found", containerName)
}

func (db *Db) InsertContainer(container models.Container) error {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	_, err := db.conn.ExecContext(
		ctx,
		fmt.Sprintf("insert into %s(name, ip_address, exposed_port) values(?, ?, ?)", containerTableName),
		container.Name, container.IpAddress, container.ExposedPort,
	)
	return err
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
