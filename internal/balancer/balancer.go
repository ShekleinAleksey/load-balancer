package balancer

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ShekleinAleksey/load-balancer/config"
)

type Balancer interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func New(backendURLs []string, algorithm string, healthCheck config.HealthCheck) (Balancer, error) {
	pool, err := NewServerPool(backendURLs)
	if err != nil {
		return nil, fmt.Errorf("cannot create serverPool: %v", err)
	}

	interval, err := time.ParseDuration(healthCheck.Interval)
	if err != nil {
		return nil, fmt.Errorf("invalid health check interval: %v", err)
	}

	timeout, err := time.ParseDuration(healthCheck.Timeout)
	if err != nil {
		return nil, fmt.Errorf("invalid health check timeout: %v", err)
	}

	go pool.HealthCheck(interval, timeout)

	switch algorithm {
	case "round-robin":
		return NewRoundRobin(pool), nil
	case "least-connections":
		return NewLeastConnections(pool), nil
	case "random":
		return NewRandom(pool), nil
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", algorithm)
	}
}
