package proxy

import (
	"net/http"
	"net/url"
	"time"
)

// Proxy is a simple http proxy entity
type Proxy interface {
	// ServeHTTP proxies incoming requests
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	// IsAvailable returns whether the proxy origin was successfully connected at the last check time.
	IsAvailable() bool
	// SetHealthCheckFunc sets the passed check func as the algorithm of checking the origin availability
	SetHealthCheckFunc(check func(addr *url.URL) bool, period time.Duration)
	// GetLoad returns the number of requests being served by the proxy at the moment
	GetLoad() int32
}
