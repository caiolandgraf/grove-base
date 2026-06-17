package ratelimiter

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newReq(ip string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RemoteAddr = ip + ":12345"
	return r
}

func check(t *testing.T, rl *Limiter, ip string) error {
	t.Helper()
	rec := httptest.NewRecorder()
	return rl.Check(rec, newReq(ip))
}

// TestNew_AllowsUpToLimit verifies that exactly `limit` requests per IP are allowed.
func TestNew_AllowsUpToLimit(t *testing.T) {
	const limit = 5
	rl := New(limit, 10*time.Second)

	allowed, denied := 0, 0
	for range limit + 3 {
		if check(t, rl, "192.168.1.1") == nil {
			allowed++
		} else {
			denied++
		}
	}

	if allowed != limit {
		t.Errorf("expected %d allowed, got %d", limit, allowed)
	}
	if denied != 3 {
		t.Errorf("expected 3 denied, got %d", denied)
	}
}

// TestNew_PerIP verifies that each IP has its own independent token bucket.
func TestNew_PerIP(t *testing.T) {
	const limit = 3
	rl := New(limit, 10*time.Second)

	for range limit {
		if check(t, rl, "10.0.0.1") != nil {
			t.Fatal("unexpected denial for IP A")
		}
	}
	if check(t, rl, "10.0.0.1") == nil {
		t.Error("IP A should be rate limited")
	}

	for range limit {
		if check(t, rl, "10.0.0.2") != nil {
			t.Error("IP B should NOT be rate limited — it has its own bucket")
		}
	}
}

// TestNew_RefillsOverTime verifies tokens replenish after waiting.
func TestNew_RefillsOverTime(t *testing.T) {
	rl := New(2, 1*time.Second)

	for range 2 {
		if check(t, rl, "172.16.0.1") != nil {
			t.Fatal("unexpected denial on initial burst")
		}
	}
	if check(t, rl, "172.16.0.1") == nil {
		t.Error("expected denial after burst")
	}

	time.Sleep(600 * time.Millisecond)

	if check(t, rl, "172.16.0.1") != nil {
		t.Error("expected one request allowed after refill")
	}
}

// TestNew_ReturnsHTTP429 verifies the error is a proper 429.
func TestNew_ReturnsHTTP429(t *testing.T) {
	rl := New(1, 10*time.Second)
	_ = check(t, rl, "1.2.3.4")

	err := check(t, rl, "1.2.3.4")
	if err == nil {
		t.Fatal("expected 429 error, got nil")
	}
}

func TestCheck_SetsRateLimitHeaders(t *testing.T) {
	rl := New(5, 30*time.Second)
	rec := httptest.NewRecorder()

	if err := rl.Check(rec, newReq("10.0.0.5")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := rec.Header().Get("X-RateLimit-Limit"); got != "5" {
		t.Fatalf("expected X-RateLimit-Limit=5, got %q", got)
	}
	if rec.Header().Get("X-RateLimit-Remaining") == "" {
		t.Fatal("expected X-RateLimit-Remaining header")
	}
}

func TestCheck_SetsRetryAfterOnDenial(t *testing.T) {
	rl := New(1, 10*time.Second)
	rec := httptest.NewRecorder()
	_ = rl.Check(rec, newReq("10.0.0.6"))

	rec = httptest.NewRecorder()
	if err := rl.Check(rec, newReq("10.0.0.6")); err == nil {
		t.Fatal("expected rate limit error")
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Fatal("expected Retry-After header on denial")
	}
}

func TestClientIP_UsesForwardedForFromTrustedProxy(t *testing.T) {
	rl := New(1, 10*time.Second, WithTrustedProxies(mustParseCIDR(t, "127.0.0.1/32")))

	req := newReq("127.0.0.1")
	req.Header.Set("X-Forwarded-For", "203.0.113.10, 10.0.0.1")

	rec := httptest.NewRecorder()
	if err := rl.Check(rec, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req.Header.Set("X-Forwarded-For", "203.0.113.10, 10.0.0.1")
	rec = httptest.NewRecorder()
	if err := rl.Check(rec, req); err == nil {
		t.Fatal("expected shared client IP to be rate limited")
	}
}

func TestMiddleware_Returns429(t *testing.T) {
	rl := New(1, 10*time.Second)
	handler := rl.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := newReq("198.51.100.1")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rec.Code)
	}
}

func mustParseCIDR(t *testing.T, raw string) *net.IPNet {
	t.Helper()
	_, network, err := net.ParseCIDR(raw)
	if err != nil {
		t.Fatalf("parse cidr: %v", err)
	}
	return network
}
