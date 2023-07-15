package querylog

import (
	"context"
	"errors"
	"github.com/nothingZero/dust"
	"strings"
)

type MiddlewareBuilder struct {
}

func InitMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (m *MiddlewareBuilder) Build() dust.Middleware {
	return func(next dust.Handler) dust.Handler {
		return func(ctx context.Context, qc *dust.QueryContext) *dust.QueryResult {
			if qc.Type == "INSERT" || qc.Type == "SELECT" {
				return next(ctx, qc)
			}
			q, err := qc.Builder.Build()
			if err != nil {
				return &dust.QueryResult{Err: err}
			}
			if strings.Contains(q.SQL, "WHERE") {
				return &dust.QueryResult{
					Err: errors.New("不允许执行不包含 WHERE 的 UPDATE 和 DELETE 语句"),
				}
			}
			return next(ctx, qc)

		}
	}
}
