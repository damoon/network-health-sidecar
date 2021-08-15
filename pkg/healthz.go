package health

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Sidecar struct {
	DNSInternal          string
	DNSExternal          string
	HTTPInternal         string
	HTTPInternalInsecure bool
	HTTPExternal         string

	dnsInternal  bool
	dnsExternal  bool
	httpInternal bool
	httpExternal bool
}

func (s *Sidecar) Run() error {
	go s.update()

	http.HandleFunc("/healthz", s.handler)
	return http.ListenAndServe(":8080", nil)
}

func (s *Sidecar) update() {
	for {
		s.dnsInternal = checkDNS(s.DNSInternal)
		s.dnsExternal = checkDNS(s.DNSExternal)
		s.httpInternal = checkHTTP(s.HTTPInternal, s.HTTPInternalInsecure)
		s.httpExternal = checkHTTP(s.HTTPExternal, false)

		time.Sleep(5 * time.Second)
	}
}

func (s *Sidecar) handler(w http.ResponseWriter, r *http.Request) {
	if !s.healthy() {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	output := fmt.Sprintf("dns internal: %v\ndns external: %v\nhttp internal: %v\nhttp external: %v\n", s.dnsInternal, s.dnsExternal, s.httpInternal, s.httpExternal)

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
