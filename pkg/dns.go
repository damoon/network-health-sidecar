package health

import (
	"context"
	"log"
	"net"
	"time"
)

func checkDNS(domain string) bool {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	ips, err := net.DefaultResolver.LookupIP(ctx, "ip4", domain)
	if err != nil {
		log.Printf("dns check (%s): %v", domain, err)
		return false
	}

	if len(ips) == 0 {
		log.Printf("dns check (%s): no IPs found", domain)
		return false
	}

	return true
}
