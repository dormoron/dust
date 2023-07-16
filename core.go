package dust

import (
	"context"
	"github.com/nothingZero/dust/internal/valuer"
	"github.com/nothingZero/dust/model"
)

type core struct {
	model   *model.Model
	dialect Dialect
	creator valuer.Creator
	r       model.Registry
	mils    []Middleware
}

func get[T any](ctx context.Context, sess Session, c core, qc *QueryContext) *QueryResult {
	var root Handler = func(ctx context.Context, qc *QueryContext) *QueryResult {
		return getHandler[T](ctx, sess, c, qc)
	}
	for i := len(c.mils) - 1; i >= 0; i-- {
		root = c.mils[i](root)
	}
	return root(ctx, qc)
}

func getHandler[T any](ctx context.Context, sess Session, c core, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{Err: err}
	}
	// 发起查询
	rows, err := sess.queryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return &QueryResult{Err: err}
	}
	if !rows.Next() {
		// 没有数据
		return &QueryResult{Err: ErrNoRows}
	}

	tp := new(T)
	val := c.creator(c.model, tp)
	err = val.SetColumns(rows)
	return &QueryResult{
		Err:    err,
		Result: tp,
	}
}

func exec(ctx context.Context, sess Session, c core, qc *QueryContext) *QueryResult {
	var root Handler = func(ctx context.Context, qc *QueryContext) *QueryResult {
		return execHandle(ctx, sess, c, qc)
	}
	for i := len(c.mils) - 1; i >= 0; i-- {
		root = c.mils[i](root)
	}
	return root(ctx, qc)
}

func execHandle(ctx context.Context, sess Session, c core, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{
			Result: Result{
				err: err,
			},
		}
	}
	res, err := sess.execContext(ctx, q.SQL, q.Args...)
	return &QueryResult{
		Result: Result{
			err: err,
			res: res,
		},
	}
}
