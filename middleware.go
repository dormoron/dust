package dust

import (
	"context"
	"github.com/nothingZero/dust/model"
)

type QueryContext struct {
	// 查询类型， 标记增删改查
	Type string

	// 查询本身
	Builder QueryBuilder

	Model *model.Model
}

type QueryResult struct {
	Result any
	Err    error
}

type Handler func(ctx context.Context, qc *QueryContext) *QueryResult

type Middleware func(next Handler) Handler
