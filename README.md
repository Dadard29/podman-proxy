# podman-proxy

`podman-proxy` is a proxy written in pure go, which redirect HTTP traffic to podman containers.

Rules are set using a web API, exposed by the proxy

## How does it work with docker

With docker, there is this [image](https://github.com/jwilder/nginx-proxy) from [Jason Wilder](https://github.com/jwilder).
It binds the docker socket into a containerized nginx proxy.
Then, it can detect the docker events and retrieve the environment configuration from the containers.

According to the env variables, it redirects the HTTP traffic to the good container.

## How does it work with podman

As podman does not use a socket, a containerized proxy cannot detect event or retrieve container environment.

Then, you cannot use a container to setup a proxy, since there is no conventional way to interact with podman from within.
(in fact, there is ssh, but configure ssh in container is pain in the ass).

I choose to use systemd instead, to run the proxy.

podman-proxy uses the [podman golang methods](https://github.com/containers/libpod) to access the containers infos.

## Installation

- install this go package
- setup the systemd unit file

## Usage

```shell script
systemctl start podman-proxy
systemctl status podman-proxy
```

- [Getting started]()
- [API documentation]()
 
