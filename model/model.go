package model

import (
	"github.com/nothingZero/dust/internal/errs"
	"reflect"
)

const (
	tagKeyColumn = "column"
)

type Model struct {
	TableName string

	Fields []*Field
	// 字段名到字段的映射
	FieldMap map[string]*Field
	// 列名到字段的映射
	ColumnMap map[string]*Field
}

type Field struct {
	// 字段名
	GoName string
	// 列名
	ColName string
	// 字段类型
	Type reflect.Type

	// 字段相对于结构体本身的偏移量
	Offset uintptr
}

func WithTableName(tableName string) Option {
	return func(model *Model) error {
		model.TableName = tableName
		return nil
	}
}

func WithColumnName(field string, colName string) Option {
	return func(model *Model) error {
		fd, ok := model.FieldMap[field]
		if !ok {
			return errs.NewErrUnknownField(fd.ColName)
		}
		fd.ColName = colName
		return nil
	}
}

type TableName interface {
	TableName() string
}
