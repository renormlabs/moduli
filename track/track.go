// Copyright (c) 2025 Renorm Labs. All rights reserved.

// Package track provides mutation-tracking primitives used by moduli.
//
// End-users import this package only when they need to inspect or stream the
// history; ordinary option-usage requires nothing beyond the root package.
package track

import (
	"encoding/json"
	"sync"
)

//------------------------------------------------------------
// Public types & helpers
//------------------------------------------------------------

// Change represents a single mutation event: the option name, the before-state,
// and the after-state.
type Change[T any] struct {
	Name   string
	Before T
	After  T
}

// Tracker is implemented by types that can record before/after mutations.
type Tracker[T any] interface {
	Track(name string, before, after T)
}

// Memory tracks all changes in memory and supports change hooks.
type Memory[T any] struct {
	mu      sync.Mutex
	history []Change[T]
	hooks   []func(Change[T])
}

// Track records a change and notifies all registered hooks. Safe for concurrent
// callers to use.
func (m *Memory[T]) Track(name string, before, after T) {
	m.mu.Lock()
	c := Change[T]{Name: name, Before: before, After: after}
	m.history = append(m.history, c)
	hooks := append([]func(Change[T]){}, m.hooks...) // copy for lock-free callbacks
	m.mu.Unlock()

	for _, h := range hooks {
		h(c)
	}
}

// History returns a copy of all recorded changes.
func (m *Memory[T]) History() []Change[T] {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Change[T], len(m.history))
	copy(out, m.history)
	return out
}

// RegisterHook adds an observer callback that runs after every change.
func (m *Memory[T]) RegisterHook(h func(Change[T])) {
	if h == nil {
		return
	}
	m.mu.Lock()
	m.hooks = append(m.hooks, h)
	m.mu.Unlock()
}

// JSON returns the tracked history encoded as JSON.
func (m *Memory[T]) JSON() ([]byte, error) { return json.Marshal(m.History()) }
