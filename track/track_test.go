// Copyright (c) 2025 Renorm Labs. All rights reserved.

package track

import (
	"bytes"
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
	m.Track("inc", data{0}, data{1})
	m.Track("inc", data{1}, data{2})

	Assert(t, Equal(int(fired.Load()), 2))

	hist := m.History()
	Assert(t, Equal(len(hist), 2))
	Assert(t, Equal(hist[1].After.V, 2))

	blob, err := m.JSON()
	Assert(t, Nil(err))

	var decoded []Change[data]
	_ = json.Unmarshal(blob, &decoded)
	Assert(t, Equal(len(decoded), 2))
}

func TestMemory_RegisterHookNil(t *testing.T) {
	m := &Memory[data]{}

	Assert(t, Not(Panics)(func() { m.RegisterHook(nil) }))

	m.Track("noop", data{}, data{})
	Assert(t, Equal(len(m.History()), 1))
}

type logCapture struct {
	buf bytes.Buffer
}

func (l *logCapture) Write(p []byte) (int, error) {
	return l.buf.Write(p)
}

func (l *logCapture) String() string {
	return l.buf.String()
}
