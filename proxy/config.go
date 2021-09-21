package proxy

import (
	"fmt"
	"os"
	"strconv"
)

type config struct {
	proxyHost string
	proxyPort int
	dbPath    string
}

const envProxyHost = "PROXY_HOST"
const envProxyPort = "PROXY_PORT"
const envDbPath = "DB_PATH"

func newConfigFromEnv() (config, error) {
	proxyHost := os.Getenv(envProxyHost)
	proxyPortStr := os.Getenv(envProxyPort)
	dbPath := os.Getenv(envDbPath)

	if proxyHost == "" {
		return config{}, fmt.Errorf("env variable not set: %s", envProxyHost)
	}
	if dbPath == "" {
		return config{}, fmt.Errorf("env variable not set: %s", envDbPath)
	}
	if proxyPortStr == "" {
		return config{}, fmt.Errorf("env variable not set: %s", envProxyPort)
	}
	proxyPort, err := strconv.Atoi(proxyPortStr)
	if err != nil {
		return config{}, fmt.Errorf("invalid type for env variable: %s = %s", envProxyPort, proxyPortStr)
	}

	proxyConfig := config{
		proxyHost: proxyHost,
		dbPath:    dbPath,
		proxyPort: proxyPort,
	}

	return proxyConfig, nil
}
