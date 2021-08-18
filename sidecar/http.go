package sidecar

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/http"
	"time"
)

func checkHTTP(url string, cert *x509.Certificate) bool {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Second,
		}).Dial,
		TLSHandshakeTimeout: time.Second,
	}

	if cert != nil {
		caCertPool := x509.NewCertPool()
		caCertPool.AddCert(cert)

		netTransport.TLSClientConfig = &tls.Config{
			RootCAs: caCertPool,
		}
	}

	var netClient = &http.Client{
		Timeout:   time.Second,
		Transport: netTransport,
	}

	resp, err := netClient.Get(url)
	if err != nil {
		log.Printf("http check (%s): %v", url, err)
		return false
	}
	defer resp.Body.Close()

	return true
}
