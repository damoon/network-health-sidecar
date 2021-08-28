# network health sidecar

Network health sidecar is intended to run as a sidecar in kubernetes pods.

It provides a healthcheck `:8080/healthz` for the readiness probe. \
Alternatively a exec probe can communicate via a local unix socket. 

Checks to verify DNS and HTTP run asynchron in a loop to allow for fast responses on the health endpoint.

## Usage examples

- [Http probe](example-http.yaml)
- [Exec probe](example-exec.yaml)

## help

``` bash
# go run cmd/network-health-server/main.go -h
NAME:
   network health sidecar - offloads network health checks form application

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --protocol value          Protocol to listen on (default: "tcp")
   --addr value              Address to listen on (default: ":8080")
   --dns-internal value      DNS domain to test cluster internal service lookups (default: "kubernetes.default.svc")
   --dns-external value      DNS domain to test external lookups (default: "cloudflare.com")
   --http-internal value     URL to test cluster internal http requests (default: "https://kubernetes.default.svc/healthz")
   --http-internal-ca value  CA to verify the internal http endpoint against (default: "/run/secrets/kubernetes.io/serviceaccount/ca.crt")
   --http-external value     URL to test external http requests (default: "https://cloudflare.com")
   --help, -h                show help (default: false)
```

``` bash
# go run cmd/network-health-client/main.go -h
NAME:
   network health sidecar - offloads network health server

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --protocol value  Protocol to listen on (default: "tcp")
   --addr value      Address to listen on (default: "http://127.0.0.1:8080")
   --metrics         Show metrics (default: false)
   --help, -h        show help (default: false)
```
