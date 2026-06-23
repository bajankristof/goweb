package slogcontext

import (
	"context"
	"log/slog"
)

type key struct{}

var argsKey = key{}

type handler struct {
	inner slog.Handler
}

type injectable[T any] interface {
	Context() context.Context
	WithContext(context.Context) T
}

func With(ctx context.Context, args ...any) context.Context {
	a, ok := ctx.Value(argsKey).([]any)
	if !ok {
		a = []any{}
	}

	a = append(a[:len(a):len(a)], args...)

	return context.WithValue(ctx, argsKey, a)
}

func Inject[T any](i injectable[T], args ...any) T {
	ctx := With(i.Context(), args...)
	return i.WithContext(ctx)
}

func NewHandler(h slog.Handler) slog.Handler {
	return &handler{inner: h}
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	args, ok := ctx.Value(argsKey).([]any)
	if ok {
		r.Add(args...)
	}

	return h.inner.Handle(ctx, r)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{inner: h.inner.WithAttrs(attrs)}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{inner: h.inner.WithGroup(name)}
}
