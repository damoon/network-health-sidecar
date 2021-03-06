package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
				Name:  "protocol",
				Value: "tcp",
				Usage: "Protocol to listen on",
			},
			&cli.StringFlag{
				Name:  "addr",
				Value: ":8080",
				Usage: "Address to listen on",
			},
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
		Action: runServer,
	}

	err := app.Run(os.Args)
	if err != nil {
		return err
	}

	return nil
}

func runServer(c *cli.Context) error {
	sidecar, err := healthChecker(c)
	if err != nil {
		return fmt.Errorf("setup health checker: %v", err)
	}

	stop := sidecar.Start()
	defer stop()

	svc := httpServer(sidecar.Mux())

	addr := c.String("addr")
	protocol := c.String("protocol")

	go mustListenAndServe(svc, protocol, addr)

	log.Println("running")

	awaitShutdown()

	log.Println("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = shutdown(ctx, svc)
	if err != nil {
		return fmt.Errorf("shutdown server: %v", err)
	}

	return nil
}

func healthChecker(c *cli.Context) (*health.Sidecar, error) {
	sidecar := &health.Sidecar{
		DNSInternal:  c.String("dns-internal"),
		DNSExternal:  c.String("dns-external"),
		HTTPInternal: c.String("http-internal"),
		HTTPExternal: c.String("http-external"),
	}

	caFile := c.String("http-internal-ca")
	if caFile != "" {
		certPEM, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("reading ca file: %v", err)
		}

		block, _ := pem.Decode([]byte(certPEM))
		if block == nil {
			return nil, fmt.Errorf("decode pem certificate: %v", err)
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse ca file: %v", err)
		}

		sidecar.HTTPInternalCA = cert
	}

	return sidecar, nil
}

func httpServer(h http.Handler) *http.Server {
	httpServer := &http.Server{
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}
	httpServer.Handler = h

	return httpServer
}

func mustListenAndServe(srv *http.Server, protocol, addr string) {
	log.Printf("starting server on %s://%s", protocol, addr)

	listener, err := net.Listen(protocol, addr)
	if err != nil {
		log.Fatal(err)
	}

	err = srv.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func awaitShutdown() {
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

func shutdown(ctx context.Context, srv *http.Server) error {
	err := srv.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
