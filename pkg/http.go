package health

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"
)

func checkHTTP(url string, insecure bool) bool {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Second,
		}).Dial,
		TLSHandshakeTimeout: time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure,
		},
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
