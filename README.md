# network health sidecar

Network health sidecar is intended to run as a sidecar in kubernetes pods.

It provides a healthcheck `:8080/healthz` for the readiness probe.

Checks to verify DNS and HTTP run asynchron in a loop to allow for fast responses on the health endpoint.

## help

``` bash
go run main.go -h
NAME:
   network health sidecar - offloads network health checks form application

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dns-internal value      DNS domain to test cluster internal service lookups (default: "kubernetes.default.svc")
   --dns-external value      DNS domain to test external lookups (default: "cloudflare.com")
   --http-internal value     URL to test cluster internal http requests (default: "https://kubernetes.default.svc/healthz")
   --http-internal-insecure  Skip https validation (default: true)
   --http-external value     URL to test external http requests (default: "https://cloudflare.com")
   --help, -h                show help (default: false)
```
