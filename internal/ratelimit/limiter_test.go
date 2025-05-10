package ratelimit

import (
	"testing"
	"time"
)

func TestTokenBucket(t *testing.T) {
	tb := NewTokenBucket(5, 0)

	// Первые 5 запросов должны проходить
	for i := 0; i < 5; i++ {
		if !tb.Allow() {
			t.Errorf("Request %d should be allowed", i)
		}
	}

	// Следующие должны блокироваться
	if tb.Allow() {
		t.Error("Request should be blocked after burst")
	}

	// Ждем пополнения
	time.Sleep(100 * time.Millisecond)
	if !tb.Allow() {
		t.Error("Request should be allowed after refill")
	}
}
