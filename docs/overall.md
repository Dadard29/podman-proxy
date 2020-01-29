# Getting started

## How does it works
The podman-proxy redirect HTTP traffic to podman containers.

As the service run with the root user, podman containers must be managed by root.

![schema]()

## Example

Be sure to be root
```shell script
root@host# whoami
root
```

Start two podman container
```shell script
root@host# podman pull docker://docker.io/jwilder/whoami:latest
root@host# podman run -d --expose 8000 --name server_1 docker.io/jwilder/whoami:latest
cbb93902d86f
root@host# podman run -d --expose 8000 --name server_2 docker.io/jwilder/whoami:latest
69a74f9cdf5e
```

Both of the containers should be reachable through their IP. Test for the container `server_1`:
```shell script
root@host# podman inspect --format '{{.NetworkSettings.IPAddress}}' server_1
10.88.0.2
root@host# curl 10.88.0.2:8000
I'm cbb93902d86f
```

The containers IDs and IPs may vary.

Now we want our containers to be reachable through the external network and with a specifc hostname, without publishing any port on the host.
Easy.

setup auth

request api

request container


