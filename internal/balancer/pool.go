package balancer

import (
	"fmt"
	"sync"
)

type ServerPool struct {
	backends []*Backend
	current  uint64
	mu       sync.RWMutex
}

// NewServerPool создает новый пул бэкендов
func NewServerPool(backendURLs []string) (*ServerPool, error) {
	pool := &ServerPool{
		backends: make([]*Backend, 0), // Инициализируем пустой слайс
		current:  0,
	}

	if len(backendURLs) == 0 {
		return nil, fmt.Errorf("backend URLs are not specified")
	}

	for _, rawURL := range backendURLs {
		backend, err := NewBackend(rawURL)
		if err != nil {
			return nil, fmt.Errorf("could not create backend: %v", err)
		}
		pool.AddBackend(backend)
	}

	return pool, nil
}

func (p *ServerPool) AddBackend(backend *Backend) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.backends = append(p.backends, backend)
}

func (p *ServerPool) NextIndex() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return int(p.current % uint64(len(p.backends)))
}

func (p *ServerPool) MarkBackendStatus(backend *Backend, alive bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	backend.SetAlive(alive)
}
