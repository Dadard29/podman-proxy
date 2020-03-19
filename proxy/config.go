package proxy

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	proxyHostKey   = "PODMAN_PROXY_HOST"
	proxyPortKey   = "PODMAN_PROXY_PORT"
	proxySecretKey = "PODMAN_PROXY_SECRET"
	proxyHostWhitelistKey = "PODMAN_PROXY_HOST_WHITELIST"
)

type Config struct {
	// the host which will not be redirected to podman container
	// the api will be available through this host
	ProxyHost  string `json:"proxy_host"`
	ProxyPort  int    `json:"proxy_port"`
	ProxyToken string `json:"proxy_token"`
	ProxyHostWhitelist []string `json:"proxy_host_whitelist"`
}

func (c Config) getAddr() string {
	if c.ProxyPort == 80 {
		return c.ProxyHost
	}
	return fmt.Sprintf("%s:%d", c.ProxyHost, c.ProxyPort)
}

func (c Config) generateToken(secret string) string {
	hash := sha256.New()
	hash.Write([]byte(secret))

	return hex.EncodeToString(hash.Sum(nil))
}

func getDefaultConfig() Config {
	defaultProxyHost, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}

	defaultProxyPort := 8080

	defaultProxyToken := "default-token"

	defaultHostWhitelist := make([]string, 0)

	return Config{
		ProxyHost:  defaultProxyHost,
		ProxyPort:  defaultProxyPort,
		ProxyToken: defaultProxyToken, // MUST be fulfilled by the user
		ProxyHostWhitelist: defaultHostWhitelist,
	}
}

func RetrieveConfig() Config {
	c := getDefaultConfig()

	if proxyHost := os.Getenv(proxyHostKey); proxyHost != "" {
		c.ProxyHost = proxyHost
	}

	if proxyPort := os.Getenv(proxyPortKey); proxyPort != "" {
		proxyPortInt, err := strconv.Atoi(proxyPort)
		if err != nil {
			log.Fatalln(err)
		}

		c.ProxyPort = proxyPortInt
	}

	if proxySecret := os.Getenv(proxySecretKey); proxySecret != "" {
		c.ProxyToken = c.generateToken(proxySecret)
		log.Println(fmt.Sprintf("token generated"))
	} else {
		log.Fatalln(
			fmt.Sprintf("You need to specify a secret into the environment variable %s !", proxySecretKey))
	}

	if proxyHostWhitelistRaw := os.Getenv(proxyHostWhitelistKey); proxyHostWhitelistRaw != "" {
		proxyHostWhitelist := strings.Split(proxyHostWhitelistRaw, ",")
		for _, s := range proxyHostWhitelist {
			c.ProxyHostWhitelist = append(c.ProxyHostWhitelist, s)
			c.ProxyHostWhitelist = append(c.ProxyHostWhitelist, fmt.Sprintf("www.%s", s))
		}
	}

	return c
}
