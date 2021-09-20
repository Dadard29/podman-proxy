package db

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	dbPath string
	conn   *sql.DB
}

func (db *Db) Init() error {
	_, err := db.conn.Exec(
		`CREATE TABLE "containers" (
			"name"	TEXT NOT NULL UNIQUE,
			"ip_address"	TEXT,
			"exposed_port"	INTEGER NOT NULL,
			PRIMARY KEY("name")
		);
		CREATE TABLE "domain_names" (
			"name"	TEXT NOT NULL UNIQUE,
			PRIMARY KEY("name")
		);
		CREATE TABLE "rules" (
			"id"	INTEGER NOT NULL UNIQUE,
			"domain_name"	TEXT NOT NULL UNIQUE,
			"container_name"	TEXT NOT NULL,
			FOREIGN KEY("domain_name") REFERENCES "domain_names"("name"),
			PRIMARY KEY("id" AUTOINCREMENT),
			FOREIGN KEY("container_name") REFERENCES "containers"("name")
		);
		CREATE TABLE "network_logs" (
			"timestamp"	INTEGER NOT NULL,
			"request_host"	TEXT NOT NULL,
			"request_method"	TEXT NOT NULL,
			"request_path"	TEXT NOT NULL,
			"request_body"	BLOB,
			"request_args"	TEXT,
			"response_status_code"	INTEGER NOT NULL,
			"response_duration"	INTEGER NOT NULL,
			"response_body"	BLOB
		);`,
	)
	if err != nil {
		return err
	}

	return nil
}

func NewDb(dbPath string) (*Db, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// enable foreign key constraints
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	conn.ExecContext(ctx, "PRAGMA foreign_keys = ON")

	return &Db{
		dbPath: dbPath,
		conn:   conn,
	}, nil
}
