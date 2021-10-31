# PODMAN PROXY

Setup a reverse proxy for incoming HTTP network.
The requests are redirected to running podman containers.
The container resolution is done with the domain name used to perform the request.
The proxy stores mapping between containers and domain names.

The available features for v2 are:
- HTTPS
- upgrades from HTTP to HTTPS (with configurable ports)
- supports containers running in pod
- logs requests and their responses
- logs informations about infra, including CPU, memory and disk

## HOW TO

### Build and run

Get and install the binary
```
$ go install github.com/Dadard29/podman-proxy
```

Setup the environment
```
PROXY_HOST="proxy-host"
PROXY_PORT="9000"
DB_PATH="./podman-proxy.sqlite3"
DEBUG="1"
UPGRADER_PORT="9001"
JWT_KEY="key"
```

Run the binary
```
$ go env GOBIN
/root/go/bin
$ /root/go/bin/podman-proxy
```

### systemd

Use this unit template:
```
[Unit]
Description=podman-proxy service to manages HTTP access to podman containers
Documentation=http://github.com/Dadard29/podman-proxy

[Service]
Environment="PROXY_HOST=<proxy_host>"
Environment="PROXY_PORT=443"
Environment="DB_PATH=/root/podman-proxy.sqlite3"
Environment="DEBUG=1"
Environment="UPGRADER_PORT=80"
Environment="JWT_KEY=<jwt_key>"
ExecStart=/root/go/bin/podman-proxy
```
