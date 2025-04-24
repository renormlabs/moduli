# moduli

`moduli` is a tiny Go helper that lets you:

1. **Build fluent configuration APIs** with the functional‑options pattern\
   (`Compose`, `If`, `Unless`, `WithDefaults`, …).
2. **Follow every mutation** with *zero boiler‑plate*: embed\
   `moduli.Trackable[T]` and you automatically get an in‑memory log.

Tracking is opt‑in and invisible when you don’t need it.

## Quick start

```bash
go get renorm.dev/moduli
```

```go
package main

import (
    "fmt"

    "renorm.dev/moduli"
    "renorm.dev/moduli/track"
)

// Domain object ------------------------------------------------------

type Rocket struct {
    moduli.Trackable[Rocket] // ← one line → full history
    Name    string
    Stage   int
}

// Domain options -----------------------------------------------------

func WithName(n string) moduli.Option[Rocket]  { return func(r *Rocket) { r.Name = n } }
func WithStage(s int) moduli.Option[Rocket]    { return func(r *Rocket) { r.Stage = s } }

// --------------------------------------------------------------------

func main() {
    r := &Rocket{}

    // Live play‑by‑play (optional)
    r.Tracker().RegisterHook(track.ConsoleHook[Rocket]())

    moduli.Apply(r,
        moduli.Named("set name", WithName("Starship")),
        moduli.Named("configure stage", WithStage(2)),
    )

    fmt.Printf("🚀 final: %+v\n", r)

    // Post‑hoc log
    for _, ev := range r.Tracker().History() {
        fmt.Printf("%s ⇒ %#v → %#v\n", ev.Name, ev.Before, ev.After)
    }
}
```

## License

[MIT](LICENSE)
