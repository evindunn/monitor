package main

import "sync"

type lockedFloat struct {
	mu    sync.Mutex
	value float64
}

func (l *lockedFloat) Set(value float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.value = value
}

func (l *lockedFloat) Get() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.value
}
