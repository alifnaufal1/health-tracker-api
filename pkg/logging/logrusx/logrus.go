package logrusx

import (
	"context"
	"health-tracker-api/pkg/logging"

	"github.com/sirupsen/logrus"
)

type logger struct{ entry *logrus.Entry }

func New(l *logrus.Logger) logging.Logger { return &logger{entry: logrus.NewEntry(l)} }

func (l *logger) WithContext(ctx context.Context) logging.Logger {
	if id, ok := logging.RequestIDFrom(ctx); ok {
		return &logger{entry: l.entry.WithField("request_id", id)}
	}
	return l
}

func (l *logger) With(key string, value any) logging.Logger {
	return &logger{entry: l.entry.WithField(key, value)}
}

func (l *logger) Debug(msg string) { l.entry.Debug(msg) }
func (l *logger) Info(msg string)  { l.entry.Info(msg) }
func (l *logger) Warn(msg string)  { l.entry.Warn(msg) }
func (l *logger) Error(msg string) { l.entry.Error(msg) }
