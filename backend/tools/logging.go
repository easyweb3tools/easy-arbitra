package tools

import (
	"context"
	"fmt"
)

type toolLogWriterKey struct{}

func WithToolLogWriter(ctx context.Context, fn func(string)) context.Context {
	return context.WithValue(ctx, toolLogWriterKey{}, fn)
}

func LogToolf(ctx context.Context, format string, args ...any) {
	fn, ok := ctx.Value(toolLogWriterKey{}).(func(string))
	if !ok || fn == nil {
		return
	}

	fn(fmt.Sprintf(format, args...))
}
