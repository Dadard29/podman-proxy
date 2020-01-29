## Golang

Be sure you have [golang](http://golang.org) installed.

## Package 

Install the podman-proxy package :
```shell script
go get github.com/Dadard29/podman-proxy
go install github.com/Dadard29/podman-proxy
```

The `podman-proxy` binary file should be available in your GOPATH:
```shell script
file $(go env GOPATH)/bin/podman-proxy
```

## systemd

Create a unit file for systemd:
```shell script
systemctl edit podman-proxy
```

Here's an example of a valid unit file.
```
[Unit]
Description=podman-proxy service to manages HTTP access to podman containers
Documentation=http://github.com/Dadard29/podman-proxy

[Service]
Environment="PODMAN_PROXY_DB=<your_go_path>/src/github.com/Dadard29/podman-proxy/api/db/podman-proxy.db"
Environment="PODMAN_PROXY_HOST=podman-proxy-host"
Environment="PODMAN_PROXY_PORT=80"
Environment="PODMAN_PROXY_SECRET=salut"
ExecStart=<your_go_path>/bin/podman-proxy
```

- Replace the `<your_go_path>` values by the output of `go env GOPATH`.

- Set the `PODMAN_PROXY_HOST` with the name under which the proxy api will be reachable.
If not set, the api will be listening with the current machine hostname.
 
- Set the `PODMAN_PROXY_PORT` to expose the service.
If not set, the api will be exposed through the port 8080.

- Set the the `PODMAN_PROXY_SECRET` value to setup the authentication token.
This variable MUST be set, otherwise the service won't start.
This secret will be use to generate the authentication token of the api. Do not share it.

More infos and details in the awesome [Archlinux Wiki](https://wiki.archlinux.org/index.php/Systemd) about systemd configuration.
