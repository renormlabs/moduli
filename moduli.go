// Copyright (c) 2025 Renorm Labs. All rights reserved.

// Package moduli embraces the functional-options pattern and adds helpers for
// composition, conditional logic, and—via the sub-package track—mutation
// tracking.
package moduli

import (
	"renorm.dev/moduli/track"
)

// Option is a functional option that mutates a *T in place.
// Typically used with [Apply] or [New] to construct or modify values fluently.
type Option[T any] func(*T)

// Noop returns an option that does nothing. Useful for conditional logic or
// default placeholders.
func Noop[T any]() Option[T] {
	return func(*T) {}
}

// IfElse applies either the yes or no option based on the result of cond. The
// selected option is applied to the target in-place.
func IfElse[T any](cond func() bool, yes, no Option[T]) Option[T] {
	return func(t *T) {
		if cond() {
			yes(t)
		} else {
			no(t)
		}
	}
}

// If applies opt only when cond() is true, implemented via [IfElse].
func If[T any](cond func() bool, opt Option[T]) Option[T] { return IfElse(cond, opt, Noop[T]()) }

// Unless applies opt only when cond() is false, implemented via [IfElse].
func Unless[T any](cond func() bool, opt Option[T]) Option[T] { return IfElse(cond, Noop[T](), opt) }

// Compose combines multiple options together into a single option, applied left
// to right. Nil options are safely ignored.
func Compose[T any](opts ...Option[T]) Option[T] {
	if len(opts) == 0 {
		return Noop[T]()
	}
	return func(t *T) {
		for _, o := range opts {
			if o != nil {
				o(t)
			}
		}
	}
}

// WithDefaults prepends defaultOpts before userOpts, then composes them, giving
// developers the ability to accept user options and still have a nice, fluent
// API for setting defaults.
func WithDefaults[T any](opts []Option[T], defaults ...Option[T]) Option[T] {
	return Compose(append(defaults, opts...)...)
}

// Apply applies each option to the target in order. If the target supports
// mutation tracking (via [Trackable]), changes are recorded.
func Apply[T any](target *T, opts ...Option[T]) {
	if target == nil {
		return
	}

	var tr track.Tracker[T]
	if ta, ok := any(target).(trackerProvider[T]); ok {
		tr = ta.provideTracker()
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if tr == nil {
			opt(target)
			continue
		}

		before := *target
		opt(target)
		tr.Track(optionName(opt), before, *target)
	}
}

// New allocates zero T, applies each [Option] via [Apply], and returns the *T.
func New[T any](opts ...Option[T]) *T {
	var v T
	Apply(&v, opts...)
	return &v
}
