package valuer

import (
	"database/sql"
	"github.com/nothingZero/dust/internal/errs"
	"github.com/nothingZero/dust/model"
	"reflect"
)

type reflectValue struct {
	model *model.Model
	// T 的指针
	val reflect.Value
}

func (r reflectValue) Field(name string) (any, error) {
	// 检测 name 是否合法
	_, ok := r.val.Type().FieldByName(name)
	if !ok {
		return nil, errs.NewErrUnknownField(name)
	}
	return r.val.FieldByName(name).Interface(), nil
}

var _ Creator = InitReflectValue

func InitReflectValue(model *model.Model, val any) Value {
	return reflectValue{
		model: model,
		val:   reflect.ValueOf(val).Elem(),
	}
}

func (r reflectValue) SetColumns(rows *sql.Rows) error {
	// 获取列名
	cs, err := rows.Columns()
	if err != nil {
		return err
	}

	// 处理结果
	vals := make([]any, 0, len(cs))
	valElems := make([]reflect.Value, 0, len(cs))
	for _, c := range cs {
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		// 反射创建一个实例
		val := reflect.New(fd.Type)
		vals = append(vals, val.Interface())
		valElems = append(valElems, val.Elem())
	}

	err = rows.Scan(vals...)
	if err != nil {
		return err
	}

	tpValue := r.val

	for i, c := range cs {
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		tpValue.FieldByName(fd.GoName).Set(valElems[i])
	}

	return err
}
