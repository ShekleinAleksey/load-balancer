package balancer

import (
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	// Создаем моки бэкендов
	backend1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("backend1"))
	}))
	defer backend1.Close()

	backend2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("backend2"))
	}))
	defer backend2.Close()

	backend3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("backend3"))
	}))
	defer backend3.Close()

	// Инициализируем пул
	backends := []string{
		backend1.URL,
		backend2.URL,
		backend3.URL,
	}

	pool, err := NewServerPool(backends)
	if err != nil {
		t.Fatalf("Failed to create server pool: %v", err)
	}

	rr := NewRoundRobin(pool)

	// Тестируем распределение
	requestCount := make(map[string]int)
	iterations := 100

	for i := 0; i < iterations; i++ {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		rr.ServeHTTP(w, r)

		// Определяем какой бэкенд ответил
		body := w.Body.String()
		requestCount[body]++
	}

	// Проверяем равномерность распределения
	expected := iterations / len(backends)
	deviation := 2 // Допустимое отклонение

	for backend, count := range requestCount {
		if math.Abs(float64(count-expected)) > float64(deviation) {
			t.Errorf("Backend %s got %d requests, expected %d with deviation %d",
				backend, count, expected, deviation)
		}
	}
	log.Println("aaaaaaaaa")
	t.Logf("Request distribution: %v", requestCount)
}
