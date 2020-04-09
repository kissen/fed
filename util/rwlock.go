package util

import "sync"

type RWLock struct {
	rlock sync.RWMutex
	wlock sync.Mutex
}

// Lock for writing.
func (l *RWLock) Lock() {
	l.wlock.Lock()
}

// Unlock when previously locked for writing.
func (l *RWLock) Unlock() {
	l.wlock.Unlock()
}

// Lock for reading.
func (l *RWLock) RLock() {
	l.rlock.RLock()
}

// Unlock when previously locked for reading.
func (l *RWLock) RUnlock() {
	l.rlock.RUnlock()
}
