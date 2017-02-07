package vhaline

import "sync"

// addr string -> *spair

type AtomicAddrToPair struct {
	m   map[string]*spair
	tex sync.RWMutex
}

func NewAtomicAddrToPair() *AtomicAddrToPair {
	return &AtomicAddrToPair{
		m: make(map[string]*spair),
	}
}

func (m *AtomicAddrToPair) Get(key string) *spair {
	m.tex.RLock()
	defer m.tex.RUnlock()
	return m.m[key]
}

func (m *AtomicAddrToPair) Get2(key string) (*spair, bool) {
	m.tex.RLock()
	defer m.tex.RUnlock()
	v, ok := m.m[key]
	return v, ok
}

func (m *AtomicAddrToPair) Set(key string, val *spair) {
	m.tex.Lock()
	defer m.tex.Unlock()
	m.m[key] = val
}

func (m *AtomicAddrToPair) Del(key string) {
	m.tex.Lock()
	defer m.tex.Unlock()
	delete(m.m, key)
}
