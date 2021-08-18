package sidecar

import "testing"

func Test_checkDNS(t *testing.T) {
	tests := []struct {
		name   string
		domain string
		want   bool
	}{
		{
			name:   "well known domain",
			domain: "localhost",
			want:   true,
		},
		{
			name:   "none existing domain",
			domain: "random-strings-fdashdoiuf43209cjf9",
			want:   false,
		},
		{
			name:   "cloudflare homepage",
			domain: "www.cloudflare.com",
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkDNS(tt.domain); got != tt.want {
				t.Errorf("checkDNS() = %v, want %v", got, tt.want)
			}
		})
	}
}
