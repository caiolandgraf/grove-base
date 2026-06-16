package ratelimiter

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/go-fuego/fuego"
	"golang.org/x/time/rate"
)

type entry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Func is a per-IP rate limiter.
// Instantiate once (e.g. as a Controller field) and call Check on every request.
//
//	type Controller struct {
//	    rl *ratelimiter.Func
//	}
//
//	func NewController(service Service) *Controller {
//	    return &Controller{rl: ratelimiter.New(5, 50*time.Second)}
//	}
//
//	func (ctrl *Controller) MyHandler(c fuego.ContextNoBody) (*MyResponse, error) {
//	    if err := ctrl.rl.Check(c.Request()); err != nil {
//	        return nil, err // automatically returns HTTP 429
//	    }
//	    // ... handler logic
//	}
type Func struct {
	mu      sync.Mutex
	entries map[string]*entry
	limit   int
	window  time.Duration
}

// New creates a per-IP Func that allows up to `limit` events per `window` for each IP.
//
//	rl := ratelimiter.New(5, 50*time.Second) // 5 requests per IP per 50 seconds
func New(limit int, window time.Duration) *Func {
	f := &Func{
		entries: make(map[string]*entry),
		limit:   limit,
		window:  window,
	}
	go f.cleanup()
	return f
}

func (f *Func) limiterFor(ip string) *rate.Limiter {
	f.mu.Lock()
	defer f.mu.Unlock()

	e, ok := f.entries[ip]
	if !ok {
		r := rate.Limit(float64(f.limit) / f.window.Seconds())
		e = &entry{limiter: rate.NewLimiter(r, f.limit)}
		f.entries[ip] = e
	}
	e.lastSeen = time.Now()
	return e.limiter
}

// cleanup removes IPs that haven't been seen for 2× the window, preventing memory leaks.
func (f *Func) cleanup() {
	ticker := time.NewTicker(f.window * 2)
	defer ticker.Stop()
	for range ticker.C {
		f.mu.Lock()
		for ip, e := range f.entries {
			if time.Since(e.lastSeen) > f.window*2 {
				delete(f.entries, ip)
			}
		}
		f.mu.Unlock()
	}
}

// Check verifies whether the request IP is within the rate limit.
// Returns a fuego.HTTPError with status 429 if the limit is exceeded.
// Use this as a one-liner inside fuego handlers.
//
//	if err := ctrl.rl.Check(c.Request()); err != nil {
//	    return nil, err
//	}
func (f *Func) Check(r *http.Request) error {
	ip := r.RemoteAddr
	// Strip port if present (e.g. "192.168.1.1:5432" → "192.168.1.1")
	if host, _, err := splitHostPort(ip); err == nil {
		ip = host
	}

	if !f.limiterFor(ip).Allow() {
		return fuego.HTTPError{
			Status: http.StatusTooManyRequests,
			Title:  "rate limit exceeded",
			Err:    errors.New("too many requests"),
		}
	}
	return nil
}

func splitHostPort(hostport string) (host, port string, err error) {
	// minimal split to avoid importing net just for this
	for i := len(hostport) - 1; i >= 0; i-- {
		if hostport[i] == ':' {
			return hostport[:i], hostport[i+1:], nil
		}
	}
	return "", "", errors.New("no port")
}
