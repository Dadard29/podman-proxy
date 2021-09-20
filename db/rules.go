package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Dadard29/podman-proxy/models"
)

const ruleTableName = "rules"

func (db *Db) InsertRule(domainName string, containerName string) error {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	_, err := db.conn.ExecContext(
		ctx,
		fmt.Sprintf("insert into %s(domain_name, container_name) values(?, ?)", ruleTableName),
		domainName, containerName,
	)
	return err
}

func (db *Db) DeleteRuleFromDomainName(domainName string) error {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	result, err := db.conn.ExecContext(
		ctx,
		fmt.Sprintf("delete from %s where domain_name == ?", ruleTableName),
		domainName,
	)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("deletion failed: rule with domain name %s not found", domainName)
	}
	return nil
}

func (db *Db) GetRuleFromDomainName(domainName string) (models.Rule, error) {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	row := db.conn.QueryRowContext(
		ctx,
		fmt.Sprintf("select * from %s where domain_name == ?", ruleTableName),
		domainName,
	)

	if row.Err() == sql.ErrNoRows {
		return models.Rule{}, fmt.Errorf("no rule found with dn: %s", domainName)
	}

	return models.NewRule(row.Scan)
}

func (db *Db) ListRules() ([]models.Rule, error) {
	rows, err := db.conn.Query(fmt.Sprintf("select * from %s", ruleTableName))
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	out := make([]models.Rule, 0)
	for rows.Next() {
		rule, err := models.NewRule(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, rule)
	}

	return out, nil
}
