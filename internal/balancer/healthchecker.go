package balancer

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type HealthChecker struct {
	client   *http.Client
	interval time.Duration
	balancer Balancer
}

func NewHealthChecker(balancer Balancer, interval time.Duration) *HealthChecker {
	return &HealthChecker{
		client:   &http.Client{Timeout: 5 * time.Second},
		interval: interval,
		balancer: balancer,
	}
}

func (p *ServerPool) HealthCheck(interval time.Duration, timeout time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		for _, backend := range p.backends {
			go func(b *Backend) {
				client := http.Client{Timeout: timeout}
				resp, err := client.Get(b.URL.String())
				if err != nil || resp.StatusCode != http.StatusOK {
					logrus.Printf("Health check failed for %s: %v", b.URL, err)
					b.SetAlive(false)
					return
				}
				b.SetAlive(true)
			}(backend)
		}
	}
}
