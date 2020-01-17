`podman-proxy` is a proxy written in pure go, which redirect HTTP traffic to podman containers.

Rules are set using a web API, exposed by the proxy 

# Usage

Run the podman-proxy container :
```bash
podman run -d -p <host-port>:80 -e PROXY_HOST=<proxy-hostname> registry.gitlab.com/dadard29/podman-proxy:latest
```

# How does it work with docker

With docker, there is this [image](https://github.com/jwilder/nginx-proxy) from [Jason Wilder](https://github.com/jwilder).
It binds the docker socket into a containerized nginx proxy.
Then, it can detect the docker events and retrieve the environment configuration from the containers.

According to the env variables, it redirects the HTTP traffic to the good container.

# How does it work with podman

As podman does not use a socket, a containerized proxy cannot detect event or retrieve container environment.

Then, you cannot use a container to setup a proxy, since there is no conventional way to interact with podman from within.
(in fact, there is ssh, but configure ssh in container is pain in the ass).

I choose to use systemd instead, to setup the proxy.

podman-proxy uses the podman golang code to access the containers infos.






