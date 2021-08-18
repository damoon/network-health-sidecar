package sidecar

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHealthy(t *testing.T) {
	endpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}))
	defer endpoint.Close()

	sidecar := &Sidecar{
		DNSInternal:  "localhost",
		DNSExternal:  "localhost",
		HTTPInternal: endpoint.URL,
		HTTPExternal: endpoint.URL,
	}

	stop := sidecar.Start()
	defer stop()

	ts := httptest.NewServer(sidecar.Mux())
	defer ts.Close()

	// wait for first request to finish
	time.Sleep(2 * time.Second)

	resp, err := ts.Client().Get(ts.URL + "/healthz")
	if err != nil {
		t.Errorf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("/health status code = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("read body of /health request: %v", err)
	}

	unhealthyBody := `dns internal: true
dns external: true
http internal: true
http external: true
`

	if string(body) != unhealthyBody {
		t.Errorf("/health body = %v, want %v", string(body), unhealthyBody)
	}
}

func TestUnhealthy(t *testing.T) {
	endpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer endpoint.Close()

	sidecar := &Sidecar{
		DNSInternal:  "random-strings-fdashdoiuf43209cjf9",
		DNSExternal:  "random-strings-fdashdoiuf43209cjf9",
		HTTPInternal: endpoint.URL,
		HTTPExternal: endpoint.URL,
	}

	stop := sidecar.Start()
	defer stop()

	ts := httptest.NewServer(sidecar.Mux())
	defer ts.Close()

	// wait for first request to finish
	time.Sleep(2 * time.Second)

	resp, err := ts.Client().Get(ts.URL + "/healthz")
	if err != nil {
		t.Errorf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("/health status code = %v, want %v", resp.StatusCode, http.StatusServiceUnavailable)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("read body of /health request: %v", err)
	}

	unhealthyBody := `dns internal: false
dns external: false
http internal: false
http external: false
`

	if string(body) != unhealthyBody {
		t.Errorf("/health body = %v, want %v", string(body), unhealthyBody)
	}
}

func TestMetrics(t *testing.T) {
	endpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}))
	defer endpoint.Close()

	sidecar := &Sidecar{
		DNSInternal:  "localhost",
		DNSExternal:  "localhost",
		HTTPInternal: endpoint.URL,
		HTTPExternal: endpoint.URL,
	}

	stop := sidecar.Start()
	defer stop()

	ts := httptest.NewServer(sidecar.Mux())
	defer ts.Close()

	// wait for first request to finish
	time.Sleep(2 * time.Second)

	resp, err := ts.Client().Get(ts.URL + "/metrics")
	if err != nil {
		t.Errorf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("/metrics status code = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("read body of /metrics request: %v", err)
	}

	needles := []string{
		"dns_external_duration_seconds_bucket",
		"dns_external_duration_seconds_sum",
		"dns_external_duration_seconds_count",
		"dns_external_failures_total",
		"dns_internal_duration_seconds_bucket",
		"dns_internal_duration_seconds_sum",
		"dns_internal_duration_seconds_count",
		"dns_internal_failures_total",
		"http_external_duration_seconds_bucket",
		"http_external_duration_seconds_sum",
		"http_external_duration_seconds_count",
		"http_external_failures_total",
		"http_internal_duration_seconds_bucket",
		"http_internal_duration_seconds_sum",
		"http_internal_duration_seconds_count",
		"http_internal_failures_total",
	}
	for _, needle := range needles {
		t.Run(needle, func(t *testing.T) {
			if !strings.Contains(string(body), needle) {
				t.Errorf("/metrics does not contain %v", needle)
			}
		})
	}
}
