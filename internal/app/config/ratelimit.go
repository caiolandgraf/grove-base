package config

import (
	"time"

	"github.com/caiolandgraf/grove-base/internal/app/helpers/ratelimiter"
)

func RateLimitSettings() ratelimiter.Settings {
	return ratelimiter.Settings{
		Read: ratelimiter.Limit{
			Max:    Env.RateLimitLimit,
			Window: time.Duration(Env.RateLimitWindowSec) * time.Second,
		},
		Write: ratelimiter.Limit{
			Max:    Env.RateLimitWriteLimit,
			Window: time.Duration(Env.RateLimitWriteWindowSec) * time.Second,
		},
		TrustedProxies: ratelimiter.ParseTrustedProxies(Env.RateLimitTrustedProxies),
	}
}
