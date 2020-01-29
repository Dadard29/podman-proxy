# Getting started

## How does it works
The podman-proxy redirect HTTP traffic to podman containers.

As the service run with the root user, podman containers must be managed by root.

![schema]()

## Example

### Context
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

### API

Remember the `PODMAN_PROXY_HOST`, `PODMAN_PROXY_PORT` and `PODMAN_PROXY_SECRET` values you've set in the unit file ?

For this example, we gonna take:
- `PODMAN_PROXY_HOST=podman-proxy-host`
- `PODMAN_PROXY_PORT=80`
- `PODMAN_PROXY_SECRET=salut`

You can request the exposed API this way:
```shell script
root@host# curl -I podman-proxy-host:80/rules/list
HTTP/1.1 401 Unauthorized
Date: Wed, 29 Jan 2020 15:52:35 GMT

```

A `401 Unauthorized` response should be raised.
For every request you make to the API, you need to specify a value for the `Authorization` header.

setup auth
token: 7Jw6NOeRvaIbvLaeoOuHWFdJfg1Ix1dxs9GttQc855E=

Retry your call now with the token in the `Authorization` header :
```shell script
root@host# curl -I -H "Authorization: Bearer ec9c3a34e791bda21bbcb69ea0eb875857497e0d48c75771b3d1adb5073ce791" podman-proxy-host:80/rules/list
HTTP/1.1 200 OK
Date: Wed, 29 Jan 2020 15:56:11 GMT
Content-Type: application/json
Content-Length: 372

```

create rule

### Container

Finally, just request the container with the hostname you've set in the rule.
Don't forget to set the HTTP port exposed by the proxy.
 
```shell script
root@host# curl -I server-host:80
I'm cbb93902d86f
```

Beautiful !



