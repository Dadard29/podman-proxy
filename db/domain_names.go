package db

import (
	"context"
	"fmt"

	"github.com/Dadard29/podman-proxy/models"
)

const domainNameTableName = "domain_names"

func (db *Db) InsertDomainName(domainName models.DomainName) error {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	_, err := db.conn.ExecContext(
		ctx,
		fmt.Sprintf("insert into %s(name) values(?)", domainNameTableName),
		domainName.Name,
	)
	return err
}

func (db *Db) DeleteDomainName(name string) error {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	result, err := db.conn.ExecContext(
		ctx,
		fmt.Sprintf("delete from %s where name == ?", domainNameTableName),
		name,
	)

	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("deletion failed: domain name with name %s not found", name)
	}

	return nil
}

func (db *Db) GetDomainName(name string) (models.DomainName, error) {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	row := db.conn.QueryRowContext(
		ctx,
		fmt.Sprintf("select * from %s where name == ?", domainNameTableName),
		name,
	)
	return models.NewDomainName(row.Scan)
}

func (db *Db) ListDomainNames() ([]models.DomainName, error) {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	rows, err := db.conn.QueryContext(
		ctx,
		fmt.Sprintf("select * from %s", domainNameTableName),
	)
	if err != nil {
		return nil, err
	}
	out := make([]models.DomainName, 0)
	for rows.Next() {
		domainName, err := models.NewDomainName(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, domainName)
	}

	return out, nil
}
