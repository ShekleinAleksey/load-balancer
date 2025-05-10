package balancer

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthCheck(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Инициализируем пул
	pool, err := NewServerPool([]string{ts.URL})
	if err != nil {
		t.Fatalf("Failed to create server pool: %v", err)
	}

	go pool.HealthCheck(1*time.Second, 1*time.Second)

	time.Sleep(2 * time.Second)

	// Проверяем состояние бэкенда
	if !pool.backends[0].IsAlive() {
		t.Error("Health check failed for healthy backend")
	}

	// имитируем смерть бэкенда
	ts.Close()
	time.Sleep(2 * time.Second)

	if pool.backends[0].IsAlive() {
		t.Error("Health check did not detect dead backend")
	}
}
