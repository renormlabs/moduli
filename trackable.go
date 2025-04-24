package moduli

import "renorm.dev/moduli/track"

// Trackable can be embedded to enable automatic mutation tracking when an
// [Option] is applied via [Apply] or [New].
type Trackable[T any] struct {
	tracker *track.Memory[T] `json:"-"` // zero value is fine
}

// Tracker returns the internal tracker, which can be used to inspect history,
// register change hooks, or marshal to JSON.
func (t *Trackable[T]) Tracker() *track.Memory[T] { return t.ensure() }

// Change represents a before/after pair recorded by a tracker. This is an alias
// of [track.Change].
type Change[T any] = track.Change[T]

func (t *Trackable[T]) ensure() *track.Memory[T] {
	if t.tracker == nil {
		t.tracker = &track.Memory[T]{}
	}
	return t.tracker
}

type trackerProvider[T any] interface{ provideTracker() track.Tracker[T] }

// provideTracker is used internally by [Apply] to enable tracking support.
// Satisfies the trackerProvider interface.
func (t *Trackable[T]) provideTracker() track.Tracker[T] { return t.ensure() }
