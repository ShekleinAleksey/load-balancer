package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ShekleinAleksey/load-balancer/config"
	"github.com/ShekleinAleksey/load-balancer/internal/balancer"
	"github.com/ShekleinAleksey/load-balancer/internal/ratelimit"
	"github.com/sirupsen/logrus"
)

type Server struct {
	cfg      *config.Config
	balancer balancer.Balancer
	limiter  ratelimit.RateLimiter
	server   *http.Server
}

func New(cfg *config.Config) (*Server, error) {
	// Инициализируем rate limiter
	rateLimiter := ratelimit.NewTokenBucket(cfg.RateLimit.Capacity, cfg.RateLimit.RefillPerSec)
	balancer, err := balancer.New(cfg.Backends, cfg.Algorithm, cfg.HealthCheck)
	if err != nil {
		return nil, err
	}

	return &Server{
		cfg:      cfg,
		balancer: balancer,
		limiter:  rateLimiter,
	}, nil
}

func (s *Server) Start() error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		s.balancer.ServeHTTP(w, r)
	})

	s.server = &http.Server{
		Addr:    ":" + s.cfg.ListenPort,
		Handler: handler,
	}

	// Канал для получения ошибок сервера
	serverErr := make(chan error, 1)
	go func() {
		logrus.Printf("Starting server on :%s", s.cfg.ListenPort)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Канал для сигналов ОС
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем либо сигнал завершения, либо ошибку сервера
	select {
	case err := <-serverErr:
		return err
	case sig := <-quit:
		logrus.Printf("Received signal: %v. Shutting down gracefully...", sig)
	}

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Останавливаем сервер
	if err := s.Stop(ctx); err != nil {
		return err
	}

	logrus.Println("Server gracefully stopped")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
