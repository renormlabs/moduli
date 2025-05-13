// Copyright (c) 2025 Renorm Labs. All rights reserved.

package moduli_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"sync/atomic"
	"testing"

	"renorm.dev/moduli"
	"renorm.dev/moduli/track"
	. "renorm.dev/observable"
)

type rocket struct {
	moduli.Trackable[*rocket]
	Name  string
	Value int
}

func withName(n string) moduli.Option[*rocket] { return func(r *rocket) { r.Name = n } }
func inc(val int) moduli.Option[*rocket]       { return func(r *rocket) { r.Value += val } }

type logCapture struct {
	buf bytes.Buffer
}

func (l *logCapture) Write(p []byte) (int, error) {
	return l.buf.Write(p)
}

func (l *logCapture) String() string {
	return l.buf.String()
}

func TestOptionHelpers(t *testing.T) {
	r := &rocket{}
	moduli.Apply(r, moduli.Compose(withName("falcon"), inc(1)))
	Assert(t, Equal(r.Name, "falcon"))
	Assert(t, Equal(r.Value, 1))

	r2 := &rocket{}
	defaults := []moduli.Option[*rocket]{withName("default")}
	user := []moduli.Option[*rocket]{inc(5)}
	moduli.Apply(r2, moduli.WithDefaults(user, defaults...))
	Assert(t, Equal(r2.Name, "default"))
	Assert(t, Equal(r2.Value, 5))

	r3 := &rocket{}
	cond := true
	moduli.Apply(r3, moduli.IfElse(func() bool { return cond }, inc(10), inc(20)))
	Assert(t, Equal(r3.Value, 10))

	cond = false
	moduli.Apply(r3, moduli.If(func() bool { return cond }, inc(100)))   // skipped
	moduli.Apply(r3, moduli.Unless(func() bool { return cond }, inc(1))) // executed
	Assert(t, Equal(r3.Value, 11))

	r4 := &rocket{Name: "noop"}
	moduli.Apply(r4, moduli.Noop[*rocket]())
	Assert(t, Equal(r4.Name, "noop"))
}

func TestApplyVariants(t *testing.T) {
	type plain struct{ N int }
	p := &plain{}
	Assert(t, Not(Panics)(func() {
		moduli.Apply(p, func(x *plain) { x.N = 42 })
	}))
	Assert(t, Equal(p.N, 42))

	r := &rocket{}
	moduli.Apply(r, withName("starship"), inc(2), nil)
	Assert(t, Equal(r.Name, "starship"))
	Assert(t, Equal(r.Value, 2))

	var rNil *rocket
	Assert(t, Not(Panics)(func() { moduli.Apply(rNil, withName("ignored")) }))
}

func TestApplyVariants2(t *testing.T) {
	type plain struct{ N int }
	p := &plain{}
	Assert(t, Not(Panics)(func() {
		moduli.Apply(p, func(x *plain) { x.N = 42 })
	}))
	Assert(t, Equal(p.N, 42))

	r := &rocket{}
	moduli.Apply(r, moduli.Named("set name", withName("starship")), inc(2), nil)
	Assert(t, Equal(len(r.Tracker().History()), 2))
	Assert(t, Equal(r.Name, "starship"))
	Assert(t, Equal(r.Value, 2))

	var rNil *rocket
	Assert(t, Not(Panics)(func() { moduli.Apply(rNil, withName("ignored")) }))
}

func TestMemoryTrackerAndJSON(t *testing.T) {
	r := &rocket{}
	var calls atomic.Int64

	r.Tracker().RegisterHook(func(moduli.Change[*rocket]) { calls.Add(1) })
	moduli.Apply(r, withName("foo"), inc(3))

	Assert(t, Equal(int(calls.Load()), 2))
	Assert(t, Equal(len(r.Tracker().History()), 2))

	blob, err := r.Tracker().JSON()
	Assert(t, Nil(err))

	var log []moduli.Change[rocket]
	_ = json.Unmarshal(blob, &log)
	Assert(t, Equal(len(log), 2))
}

func TestNamedAndConsoleHook(t *testing.T) {
	r := &rocket{}

	var buf logCapture
	r.Trackable.Tracker().RegisterHook(moduli.ConsoleHook[*rocket](moduli.WithConsoleWriter(&buf)))

	moduli.Apply(r, moduli.Named("explicit", inc(1)), moduli.Noop[*rocket]())

	h := r.Tracker().History()
	Assert(t, Equal(h[0].Name, "explicit"))
	Assert(t, Equal(h[1].Name, "option"))
}

func TestRegisterHookNil(t *testing.T) {
	r := &rocket{}
	var fired atomic.Int64

	r.Tracker().RegisterHook(nil)

	r.Tracker().RegisterHook(func(moduli.Change[*rocket]) { fired.Add(1) })

	moduli.Apply(r, inc(9))

	Assert(t, Equal(int(fired.Load()), 1))
	Assert(t, Equal(len(r.Tracker().History()), 1))
}

func TestNew(t *testing.T) {
	ptr := moduli.New(withName("newbie"), inc(7))
	Assert(t, Not(Nil)(ptr))
	Assert(t, Equal(ptr.Name, "newbie"))
	Assert(t, Equal(ptr.Value, 7))
}

func TestComposeEmpty(t *testing.T) {
	r := &rocket{Name: "initial"}

	moduli.Apply(r, moduli.Compose[*rocket]())

	Assert(t, Equal(r.Name, "initial"))
	Assert(t, Equal(r.Value, 0))
}

type data struct{ V int }

func TestSlogHook_LogsChange(t *testing.T) {
	m := &track.Memory[data]{}

	var out logCapture
	logger := slog.New(slog.NewTextHandler(&out, &slog.HandlerOptions{Level: slog.LevelInfo}))

	m.RegisterHook(moduli.SlogHook[data](logger))

	m.Track("mutate12345")

	logged := out.String()

	Assert(t, strings.Contains(logged, "option applied"))
	Assert(t, strings.Contains(logged, "mutate12345"))
}
