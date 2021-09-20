package models

import (
	"database/sql"
	"fmt"
)

type Container struct {
	Name        string `json:"name"`
	IpAddress   string `json:"ip_address"`
	ExposedPort int    `json:"exposed_port"`
}

func (c Container) String() string {
	return fmt.Sprintf("%s: %s:%d", c.Name, c.IpAddress, c.ExposedPort)
}

func NewContainer(rowCursor *sql.Rows) (Container, error) {
	var name string
	var ipAddress string
	var exposedPort int
	err := rowCursor.Scan(&name, &ipAddress, &exposedPort)
	if err != nil {
		return Container{}, err
	}

	return Container{
		Name:        name,
		IpAddress:   ipAddress,
		ExposedPort: exposedPort,
	}, nil
}
