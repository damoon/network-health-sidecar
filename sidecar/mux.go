package sidecar

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Sidecar) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/healthz", s)
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}
