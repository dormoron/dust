package valuer

import (
	"database/sql"
	"dust/internal/errs"
	"dust/model"
	"reflect"
	"unsafe"
)

type unsafeValue struct {
	model *model.Model
	// 基准地址
	address unsafe.Pointer
}

func (u unsafeValue) Field(name string) (any, error) {
	fd, ok := u.model.FieldMap[name]
	if !ok {
		return nil, errs.NewErrUnknownField(name)
	}
	fdAddress := unsafe.Pointer(uintptr(u.address) + fd.Offset)
	val := reflect.NewAt(fd.Type, fdAddress)

	return val.Elem().Interface(), nil
}

var _ Creator = InitUnsafeValue

func InitUnsafeValue(model *model.Model, val any) Value {
	address := reflect.ValueOf(val).UnsafePointer()
	return unsafeValue{
		model:   model,
		address: address,
	}
}

func (u unsafeValue) SetColumns(rows *sql.Rows) error {
	// 获取列名
	cs, err := rows.Columns()
	if err != nil {
		return err
	}

	// 处理结果
	var vals []any

	for _, c := range cs {
		fd, ok := u.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		// 反射创建一个实例
		fdAddress := unsafe.Pointer(uintptr(u.address) + fd.Offset)
		val := reflect.NewAt(fd.Type, fdAddress)
		vals = append(vals, val.Interface())
	}
	return rows.Scan(vals...)
}
