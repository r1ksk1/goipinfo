package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	dnsServer1 = "8.8.8.8"
	dnsServer2 = "8.8.4.4"
	dnsServer3 = "1.1.1.1"
	dnsServer4 = "1.0.0.1"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ipAddr := r.RemoteAddr
		hostname, _ := lookupHostname(ipAddr)
		userAgent := r.Header.Get("User-Agent")
		language := r.Header.Get("Accept-Language")

		fmt.Fprintf(w, "Client IP address: %s\n", ipAddr)
		fmt.Fprintf(w, "Client hostname: %s\n", hostname)
		fmt.Fprintf(w, "Client user agent: %s\n", userAgent)
		fmt.Fprintf(w, "Client language: %s\n", language)
	})

	http.ListenAndServe(":8448", nil)
}

func lookupHostname(ipAddr string) (string, error) {
	servers := []string{dnsServer1, dnsServer2, dnsServer3, dnsServer4}
	for _, server := range servers {
		resolver := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Second,
				}
				return d.DialContext(ctx, "udp", net.JoinHostPort(server, "53"))
			},
		}

		addr, err := net.ResolveIPAddr("ip", ipAddr)
		if err != nil {
			return "", err
		}

		names, err := resolver.LookupAddr(context.Background(), addr.IP.String())
		if err == nil && len(names) > 0 {
			// remove trailing period from hostname
			hostname := strings.TrimSuffix(names[0], ".")
			return hostname, nil
		}
	}

	return "", fmt.Errorf("could not resolve hostname for IP address %s", ipAddr)
}
