package logging

import "context"

type Logger interface {
	WithContext(ctx context.Context) Logger
	With(key string, value any) Logger
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type requestIDKey struct{}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

func RequestIDFrom(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey{}).(string)
	return id, ok
}

type nop struct{}

func Nop() Logger                                { return nop{} }
func (n nop) WithContext(context.Context) Logger { return n }
func (n nop) With(string, any) Logger            { return n }
func (nop) Debug(string)                         {}
func (nop) Info(string)                          {}
func (nop) Warn(string)                          {}
func (nop) Error(string)                         {}
