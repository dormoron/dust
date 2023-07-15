package querylog

import (
	"context"
	"github.com/nothingZero/dust"
	"log"
	"time"
)

type MiddlewareBuilder struct {
	// 慢查询阈值
	threshold time.Duration
	logFunc   func(query string, args []any)
}

func InitMiddlewareBuilder(threshold time.Duration) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(query string, args []any) {
			log.Printf("SQL: %s, Args: %v", query, args)
		},
		threshold: threshold,
	}
}

func (m *MiddlewareBuilder) LogFunc(fn func(query string, args []any)) *MiddlewareBuilder {
	m.logFunc = fn
	return m
}

func (m *MiddlewareBuilder) Build() dust.Middleware {
	return func(next dust.Handler) dust.Handler {
		return func(ctx context.Context, qc *dust.QueryContext) *dust.QueryResult {
			startTime := time.Now()
			defer func() {
				duration := time.Since(startTime)
				// 不是慢查询
				if duration < m.threshold {
					return
				}
				q, err := qc.Builder.Build()
				if err == nil {
					m.logFunc(q.SQL, q.Args)
				}

			}()

			return next(ctx, qc)
		}
	}
}
