// Copyright (c) 2025 Renorm Labs. All rights reserved.

package track

import (
	"encoding/json"
	"sync/atomic"
	"testing"

	. "renorm.dev/observable"
)

// domain type for generic tests
type data struct{ V int }

func TestMemory_Track_History_JSON(t *testing.T) {
	m := &Memory[data]{}

	var fired atomic.Int64
	m.RegisterHook(func(Change[data]) { fired.Add(1) })

	// exercise Track twice
	m.Track("inc")
	m.Track("inc")

	Assert(t, Equal(int(fired.Load()), 2))

	hist := m.History()
	Assert(t, Equal(len(hist), 2))

	blob, err := m.JSON()
	Assert(t, Nil(err))

	var decoded []Change[data]
	_ = json.Unmarshal(blob, &decoded)
	Assert(t, Equal(len(decoded), 2))
}

func TestMemory_RegisterHookNil(t *testing.T) {
	m := &Memory[data]{}

	Assert(t, Not(Panics)(func() { m.RegisterHook(nil) }))

	m.Track("noop")
	Assert(t, Equal(len(m.History()), 1))
}
