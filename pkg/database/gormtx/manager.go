package gormtx

import (
	"context"

	"gorm.io/gorm"
)

type ctxKey struct{}

type Manager struct{ db *gorm.DB }

func NewManager(db *gorm.DB) *Manager { return &Manager{db: db} }

func (m *Manager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(ctxKey{}).(*gorm.DB); ok {
		return fn(ctx)
	}
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, ctxKey{}, tx))
	})
}

func (m *Manager) DB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(ctxKey{}).(*gorm.DB); ok {
		return tx
	}
	return m.db
}
