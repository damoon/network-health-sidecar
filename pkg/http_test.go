package health

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_checkHTTP(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}))
	defer ts.Close()

	type args struct {
		url      string
		insecure bool
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
			name: "local test server",
			args: args{
				url: ts.URL,
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
			if got := checkHTTP(tt.args.url, tt.args.insecure); got != tt.want {
				t.Errorf("checkHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}
