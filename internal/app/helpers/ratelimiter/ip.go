package ratelimiter

import (
	"net"
	"net/http"
	"strings"
)

// ParseTrustedProxies parses a comma-separated list of IPs or CIDR blocks.
func ParseTrustedProxies(raw string) []*net.IPNet {
	if raw == "" {
		return nil
	}

	var out []*net.IPNet
	for part := range strings.SplitSeq(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "/") {
			_, network, err := net.ParseCIDR(part)
			if err == nil {
				out = append(out, network)
			}
			continue
		}

		ip := net.ParseIP(part)
		if ip == nil {
			continue
		}

		bits := net.IPv4len * 8
		if ip.To4() == nil {
			bits = net.IPv6len * 8
		}
		out = append(out, &net.IPNet{IP: ip, Mask: net.CIDRMask(bits, bits)})
	}

	return out
}

func clientIP(r *http.Request, trusted []*net.IPNet) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	if len(trusted) == 0 || !ipInNetworks(host, trusted) {
		return host
	}

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		for part := range strings.SplitSeq(xff, ",") {
			ip := strings.TrimSpace(part)
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	if xri := strings.TrimSpace(r.Header.Get("X-Real-IP")); xri != "" {
		if net.ParseIP(xri) != nil {
			return xri
		}
	}

	return host
}

func ipInNetworks(host string, networks []*net.IPNet) bool {
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	for _, network := range networks {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}
