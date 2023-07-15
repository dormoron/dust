package dust

type Column struct {
	name  string
	alias string
}

func (Column) expr() {}

func (Column) assign() {}

type value struct {
	val any
}

func Col(name string) Column {
	return Column{
		name: name,
	}
}
func (c Column) As(alias string) Column {
	return Column{
		name:  c.name,
		alias: alias,
	}
}

func (c Column) Add(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opAdd,
		right: valueOf(arg),
	}
}

func (c Column) Multi(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opMulti,
		right: valueOf(arg),
	}
}

func valueOf(arg any) Expression {
	switch val := arg.(type) {
	case Expression:
		return val
	default:
		return value{val: arg}
	}
}

func (c Column) selectable() {}
