package proxy

import (
	"fmt"
	"os"
	"strconv"
)

type config struct {
	proxyHost    string
	proxyPort    int
	dbPath       string
	debug        bool
	upgraderPort int
	jwtKey       string
}

const envProxyHost = "PROXY_HOST"
const envProxyPort = "PROXY_PORT"
const envDbPath = "DB_PATH"
const envDebug = "DEBUG"
const envUpgraderPort = "UPGRADER_PORT"
const envJwtKey = "JWT_KEY"

func newConfigFromEnv() (config, error) {
	proxyHost := os.Getenv(envProxyHost)
	proxyPortStr := os.Getenv(envProxyPort)
	dbPath := os.Getenv(envDbPath)
	debug := os.Getenv(envDebug)
	upgraderPortStr := os.Getenv(envUpgraderPort)
	jwtKey := os.Getenv(envJwtKey)

	vars := map[string]string{
		envProxyHost:    proxyHost,
		envProxyPort:    proxyPortStr,
		envDbPath:       dbPath,
		envDebug:        debug,
		envUpgraderPort: upgraderPortStr,
		envJwtKey:       jwtKey,
	}

	for env, value := range vars {
		if value == "" {
			return config{}, fmt.Errorf("env variable not set: %s", env)
		}
	}

	proxyPort, err := strconv.Atoi(proxyPortStr)
	if err != nil {
		return config{}, fmt.Errorf("invalid type for env variable: %s = %s", envProxyPort, proxyPortStr)
	}

	upgraderPort, err := strconv.Atoi(upgraderPortStr)
	if err != nil {
		return config{}, fmt.Errorf("invalid type for env variable: %s = %s", envProxyPort, proxyPortStr)
	}

	proxyConfig := config{
		proxyHost:    proxyHost,
		dbPath:       dbPath,
		proxyPort:    proxyPort,
		debug:        debug == "1",
		upgraderPort: upgraderPort,
		jwtKey:       jwtKey,
	}

	return proxyConfig, nil
}
