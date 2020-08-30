package api

import (
	"errors"
	"fmt"
	"go/build"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
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
	home := build.Default.GOPATH
	var dbPath string
	if home != "" {
		dbPath = fmt.Sprintf("%s/src/github.com/Dadard29/podman-proxy/api/db/podman-proxy.db", home)
	} else {
		dbPath = os.Getenv("PODMAN_PROXY_DB")
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(fmt.Sprintf("connected to db %s", dbPath))

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
		return defaultRule, errors.New(fmt.Sprintf("a rule for the container with hostname %s already exists", containerHost))
	}

	con, err := a.GetContainerFromLibpod(containerName)
	if err != nil {
		return defaultRule, err
	}

	containerIp, err := a.GetContainerIp(con)
	if err != nil {
		return defaultRule, err
	}

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
	if rule, check := a.checkRuleExistsFromHostname(containerHost); !check {
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

	if rule.ContainerHost == "" {
		return rule, errors.New("malformed rule retrieved: container name blank")
	}

	a.connector.Delete(&rule)
	return rule, nil
}

func (a *Api) UpdateRule(containerHost string, containerName string, containerPort int) (RuleModel, error) {
	defaultRule := RuleModel{}

	rule, err := a.GetRule(containerHost)
	if err != nil {
		return rule, err
	}

	con, err := a.GetContainerFromLibpod(containerName)
	if err != nil {
		return defaultRule, err
	}

	containerIp, err := a.GetContainerIp(con)
	if err != nil {
		return defaultRule, err
	}

	rule.ContainerIp = containerIp
	rule.ContainerName = containerName
	rule.ContainerPort = containerPort

	a.connector.Save(&rule)
	return rule, nil
}

func (a *Api) ListRules() ([]RuleModel, error) {
	var rulesList []RuleModel

	a.connector.Find(&rulesList)
	if rulesList == nil {
		return rulesList, errors.New("error while querying the rules list")
	}

	return rulesList, nil
}
