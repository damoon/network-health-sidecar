package health

import (
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_checkHTTP(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}))
	defer ts.Close()

	tlsts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}))
	defer tlsts.Close()

	cert, err := x509.ParseCertificate(tlsts.TLS.Certificates[0].Certificate[0])
	if err != nil {
		log.Panic(err)
	}

	type args struct {
		url    string
		caCert *x509.Certificate
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "bad domain",
			args: args{
				url: "random-strings-fdashdoiuf43209cjf9",
			},
			want: false,
		},
		{
			name: "local http test server",
			args: args{
				url: ts.URL,
			},
			want: true,
		},
		{
			name: "local https test server",
			args: args{
				url:    tlsts.URL,
				caCert: cert,
			},
			want: true,
		},
		{
			name: "cloudflare homepage",
			args: args{
				url: "https://www.cloudflare.com/",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkHTTP(tt.args.url, tt.args.caCert); got != tt.want {
				t.Errorf("checkHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}
