package dust

type op string

const (
	opEq op = "="

	opNot op = "NOT"

	opAnd op = "AND"

	opOr op = "OR"

	opLt = "<"

	opGt = ">"

	opAdd   = "+"
	opMulti = "*"
)

func (o op) String() string {
	return string(o)
}

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

func (Predicate) expr() {}
func (value) expr()     {}

func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: valueOf(arg),
	}
}

func (c Column) Lt(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opLt,
		right: valueOf(arg),
	}
}

func (c Column) Gt(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opGt,
		right: valueOf(arg),
	}
}

func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: p,
	}
}

func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opAnd,
		right: right,
	}
}

func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opOr,
		right: right,
	}
}

func (p Predicate) Add(arg any) Predicate {
	return Predicate{
		left:  p,
		op:    opAdd,
		right: valueOf(arg),
	}
}

func (p Predicate) Multi(arg any) Predicate {
	return Predicate{
		left:  p,
		op:    opMulti,
		right: valueOf(arg),
	}
}
