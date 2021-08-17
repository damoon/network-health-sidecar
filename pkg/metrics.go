package health

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	dnsInternal = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "dns_internal_duration_seconds",
		Help:    "Internal DNS checks",
		Buckets: []float64{0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
	})
	dnsInternalFailures = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dns_internal_failures_total",
		Help: "Internal failed DNS checks",
	})

	dnsExternal = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "dns_external_duration_seconds",
		Help:    "External DNS checks",
		Buckets: []float64{0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
	})
	dnsExternalFailures = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dns_external_failures_total",
		Help: "External failed DNS checks",
	})

	httpInternal = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_internal_duration_seconds",
		Help:    "Internal HTTP checks",
		Buckets: []float64{0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
	})
	httpInternalFailures = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_internal_failures_total",
		Help: "Internal failed HTTP checks",
	})

	httpExternal = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_external_duration_seconds",
		Help:    "External HTTP checks",
		Buckets: []float64{0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
	})
	httpExternalFailures = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_external_failures_total",
		Help: "External failed HTTP checks",
	})
)

func (s *Sidecar) checkDNSInternal() {
	t := time.Now()
	healthy := checkDNS(s.DNSInternal)
	d := time.Now().Sub(t).Seconds()

	dnsInternal.Observe(d)

	if !healthy || d >= 1 {
		dnsInternalFailures.Inc()
	}

	s.dnsInternal = healthy
}

func (s *Sidecar) checkDNSExternal() {
	t := time.Now()
	healthy := checkDNS(s.DNSExternal)
	d := time.Now().Sub(t).Seconds()

	dnsExternal.Observe(d)

	if !healthy || d >= 1 {
		dnsExternalFailures.Inc()
	}

	s.dnsExternal = healthy
}

func (s *Sidecar) checkHTTPInternal() {
	t := time.Now()
	healthy := checkHTTP(s.HTTPInternal, s.HTTPInternalCA)
	d := time.Now().Sub(t).Seconds()

	httpInternal.Observe(d)

	if !healthy || d >= 1 {
		httpInternalFailures.Inc()
	}

	s.httpInternal = healthy
}

func (s *Sidecar) checkHTTPExternal() {
	t := time.Now()
	healthy := checkHTTP(s.HTTPExternal, nil)
	d := time.Now().Sub(t).Seconds()

	httpExternal.Observe(d)

	if !healthy || d >= 1 {
		httpExternalFailures.Inc()
	}

	s.httpExternal = healthy
}
