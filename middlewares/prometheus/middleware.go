package prometheus

import (
	"context"
	"github.com/nothingZero/dust"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type MiddlewareBulider struct{}

type MiddlewareBuilder struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
}

func (m MiddlewareBuilder) Build() dust.Middleware {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      m.Name,
		Help:      m.Help,
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"type", "table"})
	return func(next dust.Handler) dust.Handler {
		return func(ctx context.Context, qc *dust.QueryContext) *dust.QueryResult {

			startTime := time.Now()
			defer func() {
				// 执行时间
				vector.WithLabelValues(qc.Type, qc.Model.TableName).Observe(float64(time.Since(startTime).Milliseconds()))
			}()

			return next(ctx, qc)
		}
	}
}
