package balancer

import (
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
)

type leastConnections struct {
	pool       *ServerPool
	connCounts map[*Backend]int
	mu         sync.Mutex
}

func NewLeastConnections(pool *ServerPool) *leastConnections {
	lc := &leastConnections{
		pool:       pool,
		connCounts: make(map[*Backend]int),
	}

	// Инициализация счетчиков
	for _, b := range pool.backends {
		lc.connCounts[b] = 0
	}

	return lc
}

func (lc *leastConnections) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	var selected *Backend
	minConns := int(^uint(0) >> 1) // Max int

	// Ищем бэкенд с минимальным числом соединений
	for backend := range lc.connCounts {
		if !backend.IsAlive() {
			continue
		}

		if count := lc.connCounts[backend]; count < minConns {
			minConns = count
			selected = backend
		}
	}

	if selected == nil {
		http.Error(w, "No available backends", http.StatusServiceUnavailable)
		return
	}

	// Увеличиваем счетчик
	lc.connCounts[selected]++
	defer func() { lc.connCounts[selected]-- }()
	logrus.Printf("Least-connections algorithm: routing to %s", selected.URL.Host)
	selected.Proxy.ServeHTTP(w, r)
}
