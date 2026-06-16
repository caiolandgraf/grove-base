package ratelimiter

import (
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

// TestNew_AllowsUpToLimit verifies that exactly `limit` requests per IP are allowed.
func TestNew_AllowsUpToLimit(t *testing.T) {
	const limit = 5
	rl := New(limit, 10*time.Second)

	allowed, denied := 0, 0
	for range limit + 3 {
		if rl.Check(newReq("192.168.1.1")) == nil {
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

	// Exhaust IP A
	for range limit {
		if rl.Check(newReq("10.0.0.1")) != nil {
			t.Fatal("unexpected denial for IP A")
		}
	}
	if rl.Check(newReq("10.0.0.1")) == nil {
		t.Error("IP A should be rate limited")
	}

	// IP B must still have its own full bucket
	for range limit {
		if rl.Check(newReq("10.0.0.2")) != nil {
			t.Error("IP B should NOT be rate limited — it has its own bucket")
		}
	}
}

// TestNew_RefillsOverTime verifies tokens replenish after waiting.
func TestNew_RefillsOverTime(t *testing.T) {
	rl := New(2, 1*time.Second)

	for range 2 {
		if rl.Check(newReq("172.16.0.1")) != nil {
			t.Fatal("unexpected denial on initial burst")
		}
	}
	if rl.Check(newReq("172.16.0.1")) == nil {
		t.Error("expected denial after burst")
	}

	time.Sleep(600 * time.Millisecond)

	if rl.Check(newReq("172.16.0.1")) != nil {
		t.Error("expected one request allowed after refill")
	}
}

// TestNew_ReturnsHTTP429 verifies the error is a proper 429.
func TestNew_ReturnsHTTP429(t *testing.T) {
	rl := New(1, 10*time.Second)
	_ = rl.Check(newReq("1.2.3.4"))

	err := rl.Check(newReq("1.2.3.4"))
	if err == nil {
		t.Fatal("expected 429 error, got nil")
	}
	t.Logf("error returned: %v", err)
}
