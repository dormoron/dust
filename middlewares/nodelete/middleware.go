package querylog

import (
	"context"
	"errors"
	"github.com/nothingZero/dust"
)

type MiddlewareBuilder struct {
}

func InitMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (m *MiddlewareBuilder) Build() dust.Middleware {
	return func(next dust.Handler) dust.Handler {
		return func(ctx context.Context, qc *dust.QueryContext) *dust.QueryResult {
			// 禁用 DELETE 语句
			if qc.Type == "DELETE" {
				return &dust.QueryResult{Err: errors.New("禁用 DELETE 语句")}
			}
			return next(ctx, qc)
		}
	}
}
