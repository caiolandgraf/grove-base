package ratelimiter

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-fuego/fuego"
	"golang.org/x/time/rate"
)

type entry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Limiter is a per-IP rate limiter.
// Instantiate once (e.g. as a Controller field) and call Check on every request.
//
//	type Controller struct {
//	    rl *ratelimiter.Limiter
//	}
//
//	func NewController(service Service) *Controller {
//	    return &Controller{rl: ratelimiter.New(5, 50*time.Second)}
//	}
//
//	func (ctrl *Controller) MyHandler(c fuego.ContextNoBody) (*MyResponse, error) {
//	    if err := ctrl.rl.Check(c.Response(), c.Request()); err != nil {
//	        return nil, err // automatically returns HTTP 429
//	    }
//	    // ... handler logic
//	}
type Limiter struct {
	mu             sync.Mutex
	entries        map[string]*entry
	limit          int
	window         time.Duration
	trustedProxies []*net.IPNet
}

// Option configures a Limiter.
type Option func(*Limiter)

// WithTrustedProxies enables reading the client IP from X-Forwarded-For / X-Real-IP
// when the direct peer matches one of the trusted proxies.
func WithTrustedProxies(proxies ...*net.IPNet) Option {
	return func(l *Limiter) {
		l.trustedProxies = proxies
	}
}

// New creates a per-IP Limiter that allows up to `limit` events per `window` for each IP.
//
//	rl := ratelimiter.New(5, 50*time.Second) // 5 requests per IP per 50 seconds
func New(limit int, window time.Duration, opts ...Option) *Limiter {
	if limit < 1 {
		limit = 1
	}
	if window < time.Second {
		window = time.Second
	}

	l := &Limiter{
		entries: make(map[string]*entry),
		limit:   limit,
		window:  window,
	}
	for _, opt := range opts {
		opt(l)
	}
	go l.cleanup()
	return l
}

func (l *Limiter) limiterFor(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	e, ok := l.entries[ip]
	if !ok {
		r := rate.Limit(float64(l.limit) / l.window.Seconds())
		e = &entry{limiter: rate.NewLimiter(r, l.limit)}
		l.entries[ip] = e
	}
	e.lastSeen = time.Now()
	return e.limiter
}

func (l *Limiter) cleanup() {
	ticker := time.NewTicker(l.window * 2)
	defer ticker.Stop()
	for range ticker.C {
		l.mu.Lock()
		for ip, e := range l.entries {
			if time.Since(e.lastSeen) > l.window*2 {
				delete(l.entries, ip)
			}
		}
		l.mu.Unlock()
	}
}

func (l *Limiter) writeHeaders(w http.ResponseWriter, remaining int, retryAfter time.Duration) {
	if w == nil {
		return
	}

	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(l.limit))
	if remaining < 0 {
		remaining = 0
	}
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

	if retryAfter > 0 {
		seconds := int(retryAfter.Round(time.Second).Seconds())
		if seconds < 1 {
			seconds = 1
		}
		w.Header().Set("Retry-After", strconv.Itoa(seconds))
	}
}

// Check verifies whether the request IP is within the rate limit.
// Returns a fuego.HTTPError with status 429 if the limit is exceeded.
// Pass the response writer to emit X-RateLimit-* and Retry-After headers.
//
//	if err := ctrl.rl.Check(c.Response(), c.Request()); err != nil {
//	    return nil, err
//	}
func (l *Limiter) Check(w http.ResponseWriter, r *http.Request) error {
	ip := clientIP(r, l.trustedProxies)
	lim := l.limiterFor(ip)

	reservation := lim.Reserve()
	if !reservation.OK() {
		l.writeHeaders(w, 0, l.window)
		return l.tooManyRequests()
	}

	if delay := reservation.DelayFrom(time.Now()); delay > 0 {
		reservation.Cancel()
		l.writeHeaders(w, 0, delay)
		return l.tooManyRequests()
	}

	remaining := int(lim.Tokens())
	l.writeHeaders(w, remaining, 0)
	return nil
}

func (l *Limiter) tooManyRequests() error {
	return fuego.HTTPError{
		Status: http.StatusTooManyRequests,
		Title:  "rate limit exceeded",
		Err:    errors.New("too many requests"),
	}
}

// Func is an alias kept for backward compatibility.
type Func = Limiter
