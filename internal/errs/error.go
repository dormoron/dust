package errs

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly      = errors.New("orm: 只支持指向结构体的一级指针")
	ErrNoRows           = errors.New("orm: 没有数据")
	ErrInsertZeroRow    = errors.New("orm: 插入 0 行")
	ErrNoUpdatedColumns = errors.New("orm: 未指定更新的列")
)

func NewErrFailedToRollbackTx(bizErr error, rbErr error, panicked bool) error {
	return fmt.Errorf("orm: 事务闭包回滚失败，业务错误: %w, 回滚错误: %w, 是否 panic: %t", bizErr, rbErr, panicked)
}

func NewErrUnsupportedExpression(expr any) error {
	return fmt.Errorf("orm: 不支持的表达式 %v", expr)
}
func NewErrUnsupportedAssignableType(exp any) error {
	return fmt.Errorf("orm: 不支持的 Assignable 表达式 %v", exp)
}

func NewErrUnknownField(name string) error {
	return fmt.Errorf("orm: 非法字段 %v", name)
}

func NewErrUnknownColumn(name string) error {
	return fmt.Errorf("orm: 非法列 %v", name)
}

func NewErrInvalidTagContent(pair string) error {
	return fmt.Errorf("orm: 非法标签值 %s", pair)
}

func NewErrUnsupportedAssignable(expr any) error {
	return fmt.Errorf("orm: 不支持的赋值表达式类型 %s", expr)
}
