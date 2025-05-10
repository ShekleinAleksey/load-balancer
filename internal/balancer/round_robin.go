package balancer

import (
	"net/http"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

type roundRobin struct {
	pool    *ServerPool
	counter uint64
}

func NewRoundRobin(pool *ServerPool) *roundRobin {
	return &roundRobin{pool: pool}
}

func (rr *roundRobin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	attempts := 0
	maxAttempts := len(rr.pool.backends)

	for attempts < maxAttempts {
		idx := int(atomic.AddUint64(&rr.counter, 1) % uint64(len(rr.pool.backends)))
		backend := rr.pool.backends[idx]

		if backend.IsAlive() {
			logrus.Printf("Round-robin algorithm: routing to %s", backend.URL.Host)
			backend.Proxy.ServeHTTP(w, r)
			return
		}
		attempts++
	}

	http.Error(w, "No available backends", http.StatusServiceUnavailable)
}
