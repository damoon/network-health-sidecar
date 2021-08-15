package main

import (
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
			&cli.BoolFlag{
				Name:  "http-internal-insecure",
				Value: true,
				Usage: "Skip https validation",
			},
			&cli.StringFlag{
				Name:  "http-external",
				Value: "https://cloudflare.com",
				Usage: "URL to test external http requests",
			},
		},
		Action: func(c *cli.Context) error {
			sidecar := health.Sidecar{
				DNSInternal:          c.String("dns-internal"),
				DNSExternal:          c.String("dns-external"),
				HTTPInternal:         c.String("http-internal"),
				HTTPInternalInsecure: c.Bool("http-internal-insecure"),
				HTTPExternal:         c.String("http-external"),
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
