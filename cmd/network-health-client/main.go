package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	cli "github.com/urfave/cli/v2"
)

func main() {
	err := run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(args []string) error {
	app := &cli.App{
		Name:                 "network health sidecar",
		Usage:                "offloads network health server",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "protocol",
				Value: "tcp",
				Usage: "Protocol to listen on",
			},
			&cli.StringFlag{
				Name:  "addr",
				Value: "http://127.0.0.1:8080",
				Usage: "Address to listen on",
			},
			&cli.BoolFlag{
				Name:  "metrics",
				Usage: "Show metrics",
			},
		},
		Action: queryServer,
	}

	err := app.Run(os.Args)
	if err != nil {
		return err
	}

	return nil
}

func queryServer(c *cli.Context) error {
	path := "/healthz"
	if c.Bool("metrics") {
		path = "/metrics"
	}

	addr := c.String("addr")
	protocol := strings.ToLower(c.String("protocol"))

	httpc := &http.Client{
		Timeout: time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: time.Second,
			}).Dial,
			TLSHandshakeTimeout: time.Second,
		},
	}
	url := addr + path

	if protocol == "unix" {
		httpc = &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", addr)
				},
			},
		}

		url = "http://unix" + path
	}

	resp, err := httpc.Get(url)
	if err != nil {
		return fmt.Errorf("request network health deamon on (%s) %s: %v", protocol, addr, err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %v", err)
	}

	_, err = fmt.Print(string(b))
	if err != nil {
		return fmt.Errorf("write result to stdout: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("network health deamon reports as unhealthy: http status code %d", resp.StatusCode)
	}

	return nil
}
