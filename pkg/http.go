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
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure,
		},
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
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
