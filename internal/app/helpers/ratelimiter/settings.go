package ratelimiter

import (
	"net"
	"time"
)

// Limit defines a maximum number of requests allowed within a time window.
type Limit struct {
	Max    int
	Window time.Duration
}

// Settings groups read/write limits and trusted proxy configuration.
type Settings struct {
	Read           Limit
	Write          Limit
	TrustedProxies []*net.IPNet
}
