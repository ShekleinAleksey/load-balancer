package balancer

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL   *url.URL
	Alive bool
	mu    sync.RWMutex
	Proxy *httputil.ReverseProxy
}

func NewBackend(rawURL string) (*Backend, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	return &Backend{
		URL:   u,
		Alive: true,
		Proxy: httputil.NewSingleHostReverseProxy(u),
	}, nil
}

// проверяет жив ли бэкенд
func (b *Backend) IsAlive() (alive bool) {
	b.mu.RLock()
	alive = b.Alive
	b.mu.RUnlock()
	return alive
}

// устанавливает статус бэкенда
func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	b.Alive = alive
	b.mu.Unlock()
}
