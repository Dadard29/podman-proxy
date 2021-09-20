package models

import (
	"fmt"
)

type Rule struct {
	Id            int    `json:"id"`
	DomainName    string `json:"domain_name"`
	ContainerName string `json:"container_name"`
}

func NewRule(scan func(dest ...interface{}) error) (Rule, error) {
	var id int
	var domainName string
	var containerName string
	err := scan(&id, &domainName, &containerName)
	if err != nil {
		return Rule{}, err
	}

	return Rule{
		Id:            id,
		DomainName:    domainName,
		ContainerName: containerName,
	}, nil
}

func (r Rule) String() string {
	return fmt.Sprintf("%d: %s -> %s", r.Id, r.DomainName, r.ContainerName)
}
