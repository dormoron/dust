package dust

type Deleter[T any] struct {
	builder
	table string
	where []Predicate
	sess  Session
}

func InitDelete[T any](sess Session) *Deleter[T] {
	c := sess.getCore()
	return &Deleter[T]{
		builder: builder{
			core:   c,
			quoter: c.dialect.quoter(),
		},
		sess: sess,
	}
}

func (d *Deleter[T]) Build() (*Query, error) {
	var err error
	d.model, err = d.r.Get(new(T))
	if err != nil {
		return nil, err
	}

	d.sb.WriteString("DELETE FROM ")
	if d.table == "" {
		d.quote(d.model.TableName)
	} else {
		d.sb.WriteString(d.table)
	}
	if len(d.where) > 0 {
		d.sb.WriteString(" WHERE ")
		if err := d.buildPredicates(d.where); err != nil {
			return nil, err
		}
	}
	d.sb.WriteByte(';')
	return &Query{
		SQL:  d.sb.String(),
		Args: d.args,
	}, nil
}

func (d *Deleter[T]) Form(table string) *Deleter[T] {
	d.table = table
	return d
}

func (d *Deleter[T]) Where(predicates ...Predicate) *Deleter[T] {
	d.where = predicates
	return d
}
