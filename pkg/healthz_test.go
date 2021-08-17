package health

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

	ts := httptest.NewServer(sidecar)
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

	ts := httptest.NewServer(sidecar)
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
