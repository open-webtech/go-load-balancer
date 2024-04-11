package loadbalancer

import (
	"net/http"

	"github.com/open-webtech/go-load-balancer/iterator"
)

// NewLoadBalancer returns the LoadBalancer instance with the specified iterator
func NewLoadBalancer(iterator iterator.Iterator) *LoadBalancer {
	return &LoadBalancer{
		iter: iterator,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request,err error) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		},
	}
}

// LoadBalancer handles requests distributing them between one or more proxies
type LoadBalancer struct {
	ErrorHandler func(http.ResponseWriter, *http.Request, error)

	iter iterator.Iterator
}

// ServeHTTP handles the request by the next proxy gotten from the iterator
func (l *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, err := l.iter.Next()
	if err != nil {
		l.ErrorHandler(w, r, err)
		return
	}
	p.ServeHTTP(w, r)
}
