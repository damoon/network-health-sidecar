package main

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	health "github.com/damoon/network-health-sidecar/pkg"
	cli "github.com/urfave/cli/v2"
)

func main() {
	err := run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

func run(args []string) error {
	app := &cli.App{
		Name:                 "network health sidecar",
		Usage:                "offloads network health checks form application",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "dns-internal",
				Value: "kubernetes.default.svc",
				Usage: "DNS domain to test cluster internal service lookups",
			},
			&cli.StringFlag{
				Name:  "dns-external",
				Value: "cloudflare.com",
				Usage: "DNS domain to test external lookups",
			},
			&cli.StringFlag{
				Name:  "http-internal",
				Value: "https://kubernetes.default.svc/healthz",
				Usage: "URL to test cluster internal http requests",
			},
			&cli.StringFlag{
				Name:  "http-internal-ca",
				Value: "/run/secrets/kubernetes.io/serviceaccount/ca.crt",
				Usage: "CA to verify the internal http endpoint against",
			},
			&cli.StringFlag{
				Name:  "http-external",
				Value: "https://cloudflare.com",
				Usage: "URL to test external http requests",
			},
		},
		Action: func(c *cli.Context) error {
			sidecar := health.Sidecar{
				DNSInternal:  c.String("dns-internal"),
				DNSExternal:  c.String("dns-external"),
				HTTPInternal: c.String("http-internal"),
				HTTPExternal: c.String("http-external"),
			}

			caFile := c.String("http-internal-ca")
			if caFile != "" {
				ca, err := ioutil.ReadFile(caFile)
				if err != nil {
					return fmt.Errorf("reading ca file: %v", err)
				}

				cert, err := x509.ParseCertificate(ca)
				if err != nil {
					log.Panic(err)
				}

				sidecar.HTTPInternalCA = cert
			}

			err := sidecar.Run()
			if err != nil {
				return err
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		return err
	}

	return nil
}
