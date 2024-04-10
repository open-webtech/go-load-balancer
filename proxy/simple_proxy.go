package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/open-webtech/go-load-balancer/proxy/health"
)

// NewSimpleProxy is the SimpleProxy constructor
func NewSimpleProxy(addr *url.URL) *SimpleProxy {
	return &SimpleProxy{
		proxy:  httputil.NewSingleHostReverseProxy(addr),
		health: health.NewHealthCheck(addr),
	}
}

// SimpleProxy is a simple HTTP proxy entity
type SimpleProxy struct {
	health *health.HealthCheck
	proxy  *httputil.ReverseProxy
	load   int32
}

// ServeHTTP proxies incoming requests
func (p *SimpleProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&p.load, 1)
	defer atomic.AddInt32(&p.load, -1)
	p.proxy.ServeHTTP(w, r)
}

// IsAvailable returns whether the proxy origin was successfully connected at the last check time.
func (p *SimpleProxy) IsAvailable() bool {
	return p.health.IsAvailable()
}

// SetHealthCheckFunc sets the passed check func as the algorithm of checking the origin availability
func (p *SimpleProxy) SetHealthCheckFunc(check func(addr *url.URL) bool, period time.Duration) {
	p.health.SetCheckFunc(check, period)
}

// GetLoad returns the number of requests being served by the proxy at the moment
func (p *SimpleProxy) GetLoad() int32 {
	return atomic.LoadInt32(&p.load)
}
