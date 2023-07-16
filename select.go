package dust

import (
	"context"
)

// Selectable 标记接口
// 代表查找的列或者聚合函数
type Selectable interface {
	selectable()
}

type Selector[T any] struct {
	builder
	table string

	where   []Predicate
	columns []Selectable
	groupBy []Column
	having  []Predicate
	offset  int
	limit   int
	orderBy []OrderBy

	sess Session
}

func InitSelector[T any](sess Session) *Selector[T] {
	c := sess.getCore()
	return &Selector[T]{
		builder: builder{
			core:   c,
			quoter: c.dialect.quoter(),
		},
		sess: sess,
	}
}

func (s *Selector[T]) Form(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Offset(offset int) *Selector[T] {
	s.offset = offset
	return s
}

func (s *Selector[T]) Limit(limit int) *Selector[T] {
	s.limit = limit
	return s
}

func (s *Selector[T]) Where(predicates ...Predicate) *Selector[T] {
	s.where = predicates
	return s
}

func (s *Selector[T]) Having(predicates ...Predicate) *Selector[T] {
	s.having = predicates
	return s
}

func (s *Selector[T]) GroupBy(cols ...Column) *Selector[T] {
	s.groupBy = cols
	return s
}

func (s *Selector[T]) OrderBy(ob ...OrderBy) *Selector[T] {
	s.orderBy = ob
	return s
}

func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
}

func (s *Selector[T]) Build() (*Query, error) {
	if s.model == nil {
		var err error
		s.model, err = s.r.Get(new(T))
		if err != nil {
			return nil, err
		}
	}
	s.sb.WriteString("SELECT ")
	if err := s.buildColumns(); err != nil {
		return nil, err
	}
	s.sb.WriteString(" FROM ")
	if s.table == "" {
		s.sb.WriteByte('`')
		s.sb.WriteString(s.model.TableName)
		s.sb.WriteByte('`')
	} else {
		s.sb.WriteString(s.table)
	}
	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")
		if err := s.buildPredicates(s.where); err != nil {
			return nil, err
		}
	}
	if len(s.groupBy) > 0 {
		s.sb.WriteString(" GROUP BY ")
		for i, c := range s.groupBy {
			if i > 0 {
				s.sb.WriteByte(',')
			}
			if err := s.buildColumn(c); err != nil {
				return nil, err
			}
		}
	}
	if len(s.having) > 0 {
		s.sb.WriteString(" HAVING ")
		if err := s.buildPredicates(s.having); err != nil {
			return nil, err
		}
	}
	if s.limit > 0 {
		s.sb.WriteString(" LIMIT ?")
		s.addArg(s.limit)
	}

	if s.offset > 0 {
		s.sb.WriteString(" OFFSET ?")
		s.addArg(s.offset)
	}
	s.sb.WriteByte(';')
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildColumns() error {
	if len(s.columns) == 0 {
		// 没有指定列
		s.sb.WriteByte('*')
		return nil
	}
	for i, col := range s.columns {
		if i > 0 {
			s.sb.WriteByte(',')
		}

		switch c := col.(type) {
		case Column:
			err := s.buildColumn(c)
			if err != nil {
				return err
			}
		case Aggregate:
			// 聚合函数名
			err := s.buildAggregate(c)
			if err != nil {
				return err
			}
		case RawExpr:
			s.sb.WriteString(c.raw)
			s.addArg(c.args...)
		}

	}
	return nil
}

func (s *Selector[T]) buildOrderBy() error {
	for i, ob := range s.orderBy {
		if i > 0 {
			s.sb.WriteByte(',')
		}
		err := s.buildColumn(Column{name: ob.col, alias: ""})
		if err != nil {
			return err
		}
		s.sb.WriteByte(' ')
		s.sb.WriteString(ob.order)
	}
	return nil
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	var err error
	s.model, err = s.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	res := get[T](ctx, s.sess, s.core, &QueryContext{
		Type:    "SELECT",
		Builder: s,
		Model:   s.model,
	})
	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//q, err := s.Build()
	//if err != nil {
	//	return nil, err
	//}
	//db := s.db.db
	//// 发起查询
	//rows, err := db.QueryContext(ctx, q.SQL, q.Args...)
	//if err != nil {
	//	return nil, err
	//}
	//if !rows.Next() {
	//	// 没有数据
	//	return nil, ErrNoRows
	//}
	//
	//tp := new(T)
	//var creator valuer.Creator
	//val := creator(s.model, tp)
	//err = val.SetColumns(rows)
	return nil, nil
}
