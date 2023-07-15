package dust

import (
	"github.com/nothingZero/dust/internal/errs"
	"strings"
)

type builder struct {
	core
	sb   strings.Builder
	args []any

	quoter byte
}

func (b *builder) quote(name string) {
	b.sb.WriteByte(b.quoter)
	b.sb.WriteString(name)
	b.sb.WriteByte(b.quoter)
}

func (b *builder) buildPredicates(ps []Predicate) error {
	p := ps[0]
	for i := 1; i < len(ps); i++ {
		p = p.And(ps[i])
	}
	return b.buildExpression(p)
}

func (b *builder) buildExpression(expr Expression) error {
	switch exp := expr.(type) {
	case nil:
	case Predicate:
		// 构建left
		_, ok := exp.left.(Predicate)
		if ok {
			b.sb.WriteByte('(')
		}
		if err := b.buildExpression(exp.left); err != nil {
			return err
		}
		if ok {
			b.sb.WriteByte(')')
		}
		if exp.op != "" {
			b.sb.WriteByte(' ')
			b.sb.WriteString(exp.op.String())
			b.sb.WriteByte(' ')
		}

		_, ok = exp.right.(Predicate)
		if ok {
			b.sb.WriteByte('(')
		}
		if err := b.buildExpression(exp.right); err != nil {
			return err
		}
		if ok {
			b.sb.WriteByte(')')
		}

	case Column:
		exp.alias = ""
		return b.buildColumn(exp)
	case value:
		b.sb.WriteByte('?')
		b.addArg(exp.val)
	case RawExpr:
		b.sb.WriteString(exp.raw)
		b.addArg(exp.args...)
	case Aggregate:
		return b.buildAggregate(exp)
	default:
		return errs.NewErrUnsupportedExpression(expr)
	}
	return nil

}

func (b *builder) buildAggregate(a Aggregate) error {
	b.sb.WriteString(a.fn)
	b.sb.WriteString("(`")
	fd, ok := b.model.FieldMap[a.arg]
	if !ok {
		return errs.NewErrUnknownField(a.arg)
	}
	b.sb.WriteString(fd.ColName)
	b.sb.WriteString("`)")
	if a.alias != "" {
		b.buildAs(a.alias)
	}
	return nil
}

func (b *builder) buildColumn(c Column) error {
	fd, ok := b.model.FieldMap[c.name]
	// 字段不对，列不对
	if !ok {
		return errs.NewErrUnknownField(c.name)
	}
	b.quote(fd.ColName)
	if c.alias != "" {
		b.buildAs(c.alias)
	}
	return nil
}

func (s *builder) addArg(vals ...any) {
	if len(vals) == 0 {
		return
	}
	if s.args == nil {
		s.args = make([]any, 0, 8)
	}
	s.args = append(s.args, vals...)
}

func (b *builder) buildAs(alias string) {
	if alias != "" {
		b.sb.WriteString(" AS ")
		b.quote(alias)
	}
}
