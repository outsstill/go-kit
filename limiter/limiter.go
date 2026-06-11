package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	mu       sync.Mutex
	rate     float64
	capacity float64
	tokens   float64
	last     time.Time
}

func NewTokenBucket(ratePerSecond, capacity int) *TokenBucket {
	now := time.Now()
	return &TokenBucket{rate: float64(ratePerSecond), capacity: float64(capacity), tokens: float64(capacity), last: now}
}

func (l *TokenBucket) Allow() bool {
	return l.AllowN(1)
}

func (l *TokenBucket) AllowN(n int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	l.tokens += now.Sub(l.last).Seconds() * l.rate
	if l.tokens > l.capacity {
		l.tokens = l.capacity
	}
	l.last = now
	if l.tokens >= float64(n) {
		l.tokens -= float64(n)
		return true
	}
	return false
}

type FixedWindow struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	resetAt time.Time
	count   int
}

func NewFixedWindow(limit int, window time.Duration) *FixedWindow {
	return &FixedWindow{limit: limit, window: window, resetAt: time.Now().Add(window)}
}

func (l *FixedWindow) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if now.After(l.resetAt) {
		l.count = 0
		l.resetAt = now.Add(l.window)
	}
	if l.count >= l.limit {
		return false
	}
	l.count++
	return true
}
