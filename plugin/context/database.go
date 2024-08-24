package context

import (
	"context"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	GormKey = "gorm-database"
)

var DbMap = sync.Map{}

func DatabaseMapHandle() *gorm.DB {
	value, ok := DbMap.Load(GormKey)
	if ok {
		return value.(*gorm.DB)
	} else {
		return nil
	}
}

func DatabaseSetHandle(db *gorm.DB) {
	DbMap.Store(GormKey, db)
}

func FormContext[T any](ctx context.Context, key any) (T, bool) {
	var t T
	value := ctx.Value(key)
	if value == nil {
		return t, false
	}

	t, ok := value.(T)
	return t, ok
}

func MustFromContext[T any](ctx context.Context, key any) T {
	t, ok := FormContext[T](ctx, key)
	if !ok {
		panic(fmt.Sprintf("no value found in context for key %v", key))
	}
	return t
}

func WithContext(ctx context.Context, key any, value any) context.Context {
	return context.WithValue(ctx, key, value)
}

func DatabaseHandle(ctx context.Context) *gorm.DB {
	return MustFromContext[*gorm.DB](ctx, GormKey)
}

func RegisterDatabase(db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Set(GormKey, db)
		ctx.Next()
	}
}
