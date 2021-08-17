package health

import (
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Sidecar struct {
	DNSInternal    string
	DNSExternal    string
	HTTPInternal   string
	HTTPInternalCA *x509.Certificate
	HTTPExternal   string

	dnsInternal  bool
	dnsExternal  bool
	httpInternal bool
	httpExternal bool
}

func (s *Sidecar) Start() func() {
	quit := make(chan interface{}, 4)
	done := make(chan interface{})

	go func() {
		for {
			select {
			case <-time.After(10 * time.Second):
				s.dnsInternal = checkDNS(s.DNSInternal)
			case <-quit:
				done <- struct{}{}
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-time.After(10 * time.Second):
				s.dnsExternal = checkDNS(s.DNSExternal)
			case <-quit:
				done <- struct{}{}
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-time.After(10 * time.Second):
				s.httpInternal = checkHTTP(s.HTTPInternal, s.HTTPInternalCA)
			case <-quit:
				done <- struct{}{}
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-time.After(30 * time.Second):
				s.httpExternal = checkHTTP(s.HTTPExternal, nil)
			case <-quit:
				done <- struct{}{}
				return
			}
		}
	}()

	return func() {
		for i := 0; i < 4; i++ {
			quit <- struct{}{}
		}
		for i := 0; i < 4; i++ {
			<-done
		}
	}
}

func (s *Sidecar) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !s.healthy() {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	output := fmt.Sprintf(
		"dns internal: %v\ndns external: %v\nhttp internal: %v\nhttp external: %v\n",
		s.dnsInternal,
		s.dnsExternal,
		s.httpInternal,
		s.httpExternal,
	)

	_, err := w.Write([]byte(output))
	if err != nil {
		log.Printf("write output: %v", err)
	}
}

func (s *Sidecar) healthy() bool {
	if !s.dnsInternal {
		return false
	}

	if !s.dnsExternal {
		return false
	}

	if !s.httpInternal {
		return false
	}

	if !s.httpExternal {
		return false
	}

	return true
}
