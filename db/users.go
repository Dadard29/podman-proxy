package db

import (
	"context"
	"fmt"

	"github.com/Dadard29/podman-proxy/models"
)

const usersTableName = "users"

func (db *Db) GetUser(name string) (models.User, error) {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	row := db.conn.QueryRowContext(
		ctx,
		fmt.Sprintf("select * from %s where name == ?", usersTableName),
		name,
	)
	return models.NewUser(row.Scan)
}
