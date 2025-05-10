package balancer

import (
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type random struct {
	pool   *ServerPool
	random *rand.Rand
	mu     sync.Mutex
}

func NewRandom(pool *ServerPool) *random {
	return &random{
		pool:   pool,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *random) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mu.Lock()
	defer r.mu.Unlock()

	aliveBackends := make([]*Backend, 0, len(r.pool.backends))
	for _, b := range r.pool.backends {
		if b.IsAlive() {
			aliveBackends = append(aliveBackends, b)
		}
	}

	if len(aliveBackends) == 0 {
		http.Error(w, "No available backends", http.StatusServiceUnavailable)
		return
	}

	selected := aliveBackends[r.random.Intn(len(aliveBackends))]
	logrus.Printf("Random algorithm: routing to %s", selected.URL.Host)
	selected.Proxy.ServeHTTP(w, req)
}
