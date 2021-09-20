package proxy

type config struct {
	proxyHost string
	dbPath    string
	addr      string
}

const envProxyHost = "PROXY_HOST"
const envDbPath = "DB_PATH"
const envAddr = "ADDR"
