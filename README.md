# Foreman
It is a [foreman](https://github.com/ddollar/foreman) implementation in GO.

## Description
Foreman is a manager for [Procfile-based](https://en.wikipedia.org/wiki/Procfs) applications. Its aim is to abstract away the details of the Procfile format, and allow you to run your services directly.

## Features
- Run procfile-backed applications.
- Able to run with dependency resolution.

## Procfile
Procfile is simply `key: value` format like:
```yaml
app1:
    cmd: ping -c 1 google.com
    run_once: true
    checks:
        cmd: sleep 3
    deps: 
        - redis

app2:
    cmd: ping -c 50 yahoo.com
    checks:
        cmd: sleep 4

redis:
    cmd: redis-server --port 6010
    checks:
        cmd: redis-cli -p 6010 ping
        tcp_ports: [6010]
```
**Here** we defined three services
- `app1` service, executes command `ping -c 1 google.com`, and `run_once` to not respawn it if it ever dies, we say it depends on `redis` service, so `redis` needs to execute first, and its sleep command is to sleep 3 seconds.

- `app2` service, executes command `ping -c 1 yahoo.com`, and when it's done it will respawn, and its check command is to sleep 3 seconds.

- `redis` service, executes command `redis-server --port 6010`, and when it's done it will respawn, and its check is checking for port `6010` to be listening and also there's a `redis-cli -p 6010 ping` check command to validate if redis ok.

## How to use
**First:** modify the procfile with processes or services you want to run.

**second**: simply run with command: 
```sh
go run *.go
```