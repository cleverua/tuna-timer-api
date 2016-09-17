package utils

import (
	"context"

	"github.com/jinzhu/gorm"
)

const (
	// ContextKeyDBTransaction is a key under which a request-scoped DB transaction is stored in Context
	ContextKeyDBTransaction = "dbTransaction"
)

// GetDBTransactionFromContext returns request-scoped DB transaction from the context
func GetDBTransactionFromContext(ctx context.Context) *gorm.DB {
	return ctx.Value(ContextKeyDBTransaction).(*gorm.DB)
}

// SetDBTransactionInContext saves request-scoped DB transaction in context
func SetDBTransactionInContext(parentContext context.Context, dbTransaction *gorm.DB) context.Context {
	return context.WithValue(parentContext, ContextKeyDBTransaction, dbTransaction)
}
