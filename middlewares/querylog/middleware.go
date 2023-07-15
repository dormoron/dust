package querylog

import (
	"context"
	"github.com/nothingZero/dust"
	"log"
)

type MiddlewareBuilder struct {
	logFunc func(query string, args []any)
}

func InitMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(query string, args []any) {
			log.Printf("SQL: %s, Args: %v", query, args)
		},
	}
}

func (m *MiddlewareBuilder) LogFunc(fn func(query string, args []any)) *MiddlewareBuilder {
	m.logFunc = fn
	return m
}

func (m *MiddlewareBuilder) Build() dust.Middleware {
	return func(next dust.Handler) dust.Handler {
		return func(ctx context.Context, qc *dust.QueryContext) *dust.QueryResult {
			q, err := qc.Builder.Build()
			if err != nil {
				return &dust.QueryResult{
					Err: err,
				}
			}
			m.logFunc(q.SQL, q.Args)
			res := next(ctx, qc)

			return res
		}
	}
}
