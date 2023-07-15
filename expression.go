package dust

// Expression 标记接口
type Expression interface {
	expr()
}

// RawExpr 代表原生表达式
type RawExpr struct {
	raw  string
	args []any
}

func Raw(expr string, args ...any) RawExpr {
	return RawExpr{
		raw:  expr,
		args: args,
	}
}

func (r RawExpr) selectable() {}

func (r RawExpr) expr() {}

func (r RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}
