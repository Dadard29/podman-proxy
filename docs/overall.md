# Getting started

## How does it works
The podman-proxy redirect HTTP traffic to podman containers.

As the service run with the root user, podman containers must be managed by root.

![schema](https://raw.githubusercontent.com/Dadard29/podman-proxy/master/docs/images/podman-proxy.svg)

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
root@host# curl -i podman-proxy-host:80/rules/list
HTTP/1.1 401 Unauthorized
Date: Wed, 29 Jan 2020 15:52:35 GMT
Content-Length: 0

```

### Authentication

A `401 Unauthorized` response should be raised.
For every request you make to the API, you need to specify a Bearer token for the `Authorization` header.

To get this token, you need to use the secret, and send it to the endpoint `/auth`:
```shell script
root@host# curl -i -d `{"Secret": "<your_secret>"}` podman-proxy-host:80/auth
HTTP/1.1 200 OK
Date: Thu, 30 Jan 2020 10:47:23 GMT
Content-Length: 64
Content-Type: text/plain; charset=utf-8

<the_api_token>
```


Then, retry your call with the token in the `Authorization` header :
```shell script
root@host# curl -i -H "Authorization: Bearer <the_api_token>" podman-proxy-host:80/rules/list
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 30 Jan 2020 10:52:00 GMT
Content-Length: 372

{
  "Status": true,
  "Message": "rule list retrieved",
  "Rule": []
}
```

The authentication is working !

### Rule creation
```shell script
root@host# curl -i -d '{"ContainerName": "server_1", "ContainerHost": "server-1-host", "ContainerPort": 8000}' -H "Authorization: Bearer <your_api_token>" podman-proxy-host/rules
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 30 Jan 2020 11:07:50 GMT
Content-Length: 159

{
  "Status": true,
  "Message": "rule created",
  "Rule": {
    "ContainerHost": "server-1-host",
    "ContainerName": "server_1",
    "ContainerIp": "10.88.0.2",
    "ContainerPort": 8000
  }
}
```

You can set several rules for the same container. It means the container can be reached by different hostnames.


### Container

Finally, just request the container with the hostname you've set in the rule.
Don't forget to set the HTTP port exposed by the proxy (in this case, 80).
 
```shell script
root@host# curl server-1-host:80
I'm cbb93902d86f
```

Beautiful !



