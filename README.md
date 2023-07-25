# Secret-Tunnel

![Golang](https://img.shields.io/github/actions/workflow/status/starudream/secret-tunnel/golang.yml?label=golang&style=for-the-badge)
![Release](https://img.shields.io/github/actions/workflow/status/starudream/secret-tunnel/release.yml?label=release&style=for-the-badge)
![Release](https://img.shields.io/github/v/release/starudream/secret-tunnel?include_prereleases&sort=semver&style=for-the-badge)
![License](https://img.shields.io/github/license/starudream/secret-tunnel?style=for-the-badge)

## Usage

### Environment

| Name      | Type   | Comment                                |
|-----------|--------|----------------------------------------|
| ST_DEBUG  | bool   | (global) show verbose information      |
| ST_DB_DSN | string | (server) sqlite database file location |

### Server

```text
Usage:
  server [flags]
  server [command]

Available Commands:
  client      Manage clients
  help        Help about any command
  task        Manage tasks

Flags:
      --addr string    server address (default "0.0.0.0:9797")
      --api string     api address (default "127.0.0.1:9799")
  -h, --help           help for server
      --token string   api token
  -v, --version        version for server

Use "server [command] --help" for more information about a command.
```

### Client

```text
Usage:
  client [flags]
  client [command]

Available Commands:
  help        Help about any command
  service     Run as a service

Flags:
      --addr string   server address (default "127.0.0.1:9797")
      --dns string    dns server (default "8.8.8.8")
  -h, --help          help for client
      --key string    auth key
  -v, --version       version for client

Use "client [command] --help" for more information about a command.
```

## Docker

### Server [Docker Hub](https://hub.docker.com/r/starudream/secret-tunnel-server)

![Version](https://img.shields.io/docker/v/starudream/secret-tunnel-server?sort=semver&style=for-the-badge)
![Size](https://img.shields.io/docker/image-size/starudream/secret-tunnel-server?sort=semver&style=for-the-badge)
![Pull](https://img.shields.io/docker/pulls/starudream/secret-tunnel-server?style=for-the-badge)

### Client [Docker Hub](https://hub.docker.com/r/starudream/secret-tunnel-client)

![Version](https://img.shields.io/docker/v/starudream/secret-tunnel-client?sort=semver&style=for-the-badge)
![Size](https://img.shields.io/docker/image-size/starudream/secret-tunnel-client?sort=semver&style=for-the-badge)
![Pull](https://img.shields.io/docker/pulls/starudream/secret-tunnel-client?style=for-the-badge)

## Example

### Server

```shell
sts
```

```shell
curl \
--location \
--request POST 'http://127.0.0.1:9799/client' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "发送端"
}'
```

```shell
sts client create --name 发送端
```

```shell
curl \
--location \
--request POST 'http://127.0.0.1:9799/client' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "接收端"
}'
```

```shell
sts client create --name 接收端
```

```shell
curl \
--location \
--request POST 'http://127.0.0.1:9799/task' \
--header 'Content-Type: application/json' \
--data-raw '{
    "client_id": 1,
    "name": "ssh",
    "addr": "127.0.0.1:22"
}'
```

```shell
sts task create --client-id 1 --name ssh --addr 127.0.0.1:22
```

### Client

- 发送端

```shell
stc --addr 127.0.0.1:9797 --key fb9a318168714565993f75b97e6af907
```

- 发送端（服务）

```shell
stc service --user install --addr 127.0.0.1:9797 --key fb9a318168714565993f75b97e6af907
stc service --user start
stc service --user stop
stc service --user uninstall
```

- 接收端

```shell
stc --addr 127.0.0.1:9797 --key ef335f0c7a9643d19d06591672576f46
```

```yaml
tasks:
  - name: ssh
    address: 127.0.0.1:2222
    secret: ab6080af85f44772a0613fd20c09397d
```

### Data Transfer

```text
ssh:22 <-> stc(发送端) <-> sts(9797) <-> stc(接收端) <-> ssh:2222
```

## License

[Apache License 2.0](./LICENSE)
