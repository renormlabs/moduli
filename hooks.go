package moduli

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"renorm.dev/moduli/track"
)

// SlogHookConfig allows developers to configure the behavior of SlogHook.
type SlogHookConfig struct {
	logger *slog.Logger
	level  slog.Level
	msg    string
}

// WithSlogLogger sets the slog logger instance to be logger.
func WithSlogLogger(logger *slog.Logger) Option[SlogHookConfig] {
	return func(target *SlogHookConfig) {
		target.logger = logger
	}
}

// WithSlogLevel sets the logging level to be level.
func WithSlogLevel(level slog.Level) Option[SlogHookConfig] {
	return func(target *SlogHookConfig) {
		target.level = level
	}
}

// WithSlogMessage sets the message to be msg.
func WithSlogMessage(msg string) Option[SlogHookConfig] {
	return func(target *SlogHookConfig) {
		target.msg = msg
	}
}

// SlogHook returns a change hook that logs each mutation using [log/slog].
// You can pass a custom [slog.Logger] or nil to use the default logger.
func SlogHook[T any](logger *slog.Logger, opts ...Option[SlogHookConfig]) func(track.Change[T]) {
	cfg := New(WithDefaults(
		opts,
		WithSlogLogger(slog.Default()),
		WithSlogLevel(slog.LevelInfo),
		WithSlogMessage("moduli option applied")))
	return func(c track.Change[T]) {
		logger.LogAttrs(nil, cfg.level,
			cfg.msg,
			slog.String("name", c.Name),
			slog.Any("before", c.Before),
			slog.Any("after", c.After),
		)
	}
}

// ConsoleHookConfig allows developers to configure the behavior of ConsoleHook.
type ConsoleHookConfig struct {
	writer io.Writer
}

// WithConsoleWriter sets the output to the given [io.Writer].
func WithConsoleWriter(writer io.Writer) Option[ConsoleHookConfig] {
	return func(target *ConsoleHookConfig) {
		target.writer = writer
	}
}

// ConsoleHook returns a change hook that logs each mutation to [os.Stdout] or
// the configured [io.Writer].
func ConsoleHook[T any](opts ...Option[ConsoleHookConfig]) func(track.Change[T]) {
	cfg := New(WithDefaults(
		opts,
		WithConsoleWriter(os.Stdout),
	))
	return func(c track.Change[T]) {
		fmt.Fprintf(cfg.writer, "%s:\n\tbefore: %#v\n\tafter:  %#v\n\n", c.Name, c.Before, c.After)
	}
}
