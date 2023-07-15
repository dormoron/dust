package opentelemetry

import (
	"context"
	"fmt"
	"github.com/nothingZero/dust"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "github.com/nothingZero/dust/middlewares/opentelemetry"

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func (m *MiddlewareBuilder) Build() dust.Middleware {
	if m.Tracer == nil {
		otel.GetTracerProvider().Tracer(instrumentationName)
	}

	return func(next dust.Handler) dust.Handler {
		return func(ctx context.Context, qc *dust.QueryContext) *dust.QueryResult {

			tbl := qc.Model.TableName
			spanCtx, span := m.Tracer.Start(ctx, fmt.Sprintf("%s-%s", qc.Type, tbl))
			defer span.End()

			q, _ := qc.Builder.Build()
			if q != nil {
				span.SetAttributes(attribute.String("SQL", q.SQL))
			}
			span.SetAttributes(attribute.String("table", tbl))
			span.SetAttributes(attribute.String("component", "orm"))

			res := next(spanCtx, qc)
			if res.Err != nil {
				span.RecordError(res.Err)
			}
			return res
		}
	}
}
