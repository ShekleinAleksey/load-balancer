package ratelimit

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type RateLimiter interface {
	Allow() bool
}

type TokenBucket struct {
	capacity   int           // Максимальное количество токенов
	tokens     int           // Текущее количество токенов
	refillRate time.Duration // Интервал пополнения 1 токена
	lastRefill time.Time     // Время последнего пополнения
	mu         sync.Mutex
}

func NewTokenBucket(rateLimit int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity:   rateLimit,
		tokens:     rateLimit,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow проверяет, разрешен ли запрос
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Пополняем токены
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tokensToAdd := int(elapsed / tb.refillRate)

	if tokensToAdd > 0 {
		tb.tokens = min(tb.tokens+tokensToAdd, tb.capacity)
		tb.lastRefill = now
	}

	// Проверяем доступность токенов
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	logrus.Println("Rate limit exceeded")
	return false
}
