package slogctx

import (
	"context"
	"log/slog"
)

type contextKey struct{}

var attrsKey = contextKey{}

type handler struct {
	inner slog.Handler
}

// WithAttr adds a single key-value attribute to the context, which will be included in all log records handled by the slog.Handler created with New.
func WithAttr(ctx context.Context, key string, value any) context.Context {
	return WithAttrs(ctx, slog.Any(key, value))
}

// WithAttrs adds the provided attributes to the context, which will be included in all log records handled by the slog.Handler created with New.
func WithAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	a, ok := ctx.Value(attrsKey).([]slog.Attr)
	if !ok {
		a = []slog.Attr{}
	}

	a = append(a[:len(a):len(a)], attrs...)

	return context.WithValue(ctx, attrsKey, a)
}

// NewHandler creates a new slog.Handler that enriches log records with attributes from the context.
func NewHandler(h slog.Handler) slog.Handler {
	return &handler{inner: h}
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	attrs, ok := ctx.Value(attrsKey).([]slog.Attr)
	if ok {
		r.AddAttrs(attrs...)
	}

	return h.inner.Handle(ctx, r)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{inner: h.inner.WithAttrs(attrs)}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{inner: h.inner.WithGroup(name)}
}
