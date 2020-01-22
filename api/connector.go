package api

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

type RuleModel struct {
	ContainerHost string `gorm:"index:container_host;primary_key"`
	ContainerName string `gorm:"index:container_name"`
	ContainerIp   string `gorm:"index:container_ip"`
	ContainerPort int    `gorm:"index:container_port"`
}

func (RuleModel) TableName() string {
	return "podman-proxy-rules"
}

func newConnector() *gorm.DB {
	// FIXME
	db, err := gorm.Open("sqlite3", "/home/dadard/go/src/github.com/Dadard29/podman-proxy/api/db/podman-proxy.db")
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

func (a *Api) checkRuleExistsFromName(containerName string) (RuleModel, bool) {
	rule := RuleModel{
		ContainerName: containerName,
	}
	a.connector.First(&rule)
	if rule.ContainerIp == "" {
		return rule, false
	}

	return rule, true
}

func (a *Api) checkRuleExistsFromHostname(host string) (RuleModel, bool) {
	rule := RuleModel{
		ContainerHost: host,
	}
	a.connector.First(&rule)
	if rule.ContainerIp == "" || rule.ContainerHost != host {
		return RuleModel{}, false
	}

	return rule, true
}

func (a *Api) CreateRule(containerName string, containerPort int, containerHost string) (RuleModel, error) {
	defaultRule := RuleModel{}

	if _, check := a.checkRuleExistsFromHostname(containerHost); check {
		return defaultRule, errors.New(fmt.Sprintf("a rule for the container with name %s already exists", containerName))
	}

	con, err := a.GetContainerFromLibpod(containerName)
	if err != nil {
		return defaultRule, err
	}

	containerIp, err := a.GetContainerIp(con)
	if err != nil {
		return defaultRule, err
	}

	//portsExposed, err := con.PortMappings()
	//if err != nil {
	//	return defaultRule, err
	//}

	//found := false
	//for _, p := range portsExposed {
	//	if int (p.ContainerPort) == containerPort {
	//		found = true
	//	}
	//}
	//
	//if ! found {
	//	return defaultRule, errors.New(fmt.Sprintf("this port (%d) is not exposed by the container", containerPort))
	//}

	newRule := RuleModel{
		ContainerName: containerName,
		ContainerHost: containerHost,
		ContainerIp:   containerIp,
		ContainerPort: containerPort,
	}

	a.connector.Create(&newRule)
	return newRule, nil
}

func (a *Api) GetRule(containerHost string) (RuleModel, error) {
	if rule, check := a.checkRuleExistsFromHostname(containerHost); ! check {
		return rule, errors.New(fmt.Sprintf("a rule with hostname %s has not been found", containerHost))
	} else {
		return rule, nil
	}
}

func (a *Api) DeleteRule(containerHost string) (RuleModel, error) {
	rule, err := a.GetRule(containerHost)
	if err != nil {
		return rule, err
	}

	// check if the rule has correctly been retrieved; else, all record could be deleted:
	// https://gorm.io/docs/delete.html

	if rule.ContainerName == "" {
		return rule, errors.New("malformed rule retrieved: container name blank")
	}

	a.connector.Delete(&rule)
	return rule, nil
}
