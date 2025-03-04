package dag

import (
	"context"
	"fmt"

	"go-dag/basic"
)

// DataWrapper 原料结果包装器
type DataWrapper[T any] struct {
	Data   T
	Err    error
	Reason string
}

func DWOf[T any](data T) DataWrapper[any] {
	return DataWrapper[any]{Data: data}
}

func DWEmpty[T any]() DataWrapper[any] {
	return DataWrapper[any]{Data: basic.Zero[T]()}
}
func DWErrorOf1[T any](err error) DataWrapper[any] {
	return DataWrapper[any]{Data: basic.Zero[T](), Err: err}
}

func DWErrorOf2[T any](err error, reason string) DataWrapper[any] {
	return DataWrapper[any]{Data: basic.Zero[T](), Err: err, Reason: reason}
}

func (dw DataWrapper[T]) LogCustom(ctx context.Context, module string, msg string) DataWrapper[T] {
	if dw.Err != nil {
		fmt.Printf("[Node][%v][%T] pack error: %v", module, dw.Data, msg)
	}
	return dw
}

func (dw DataWrapper[T]) Log(ctx context.Context, module string) DataWrapper[T] {
	return dw.LogCustom(ctx, module, "")
}

func (dw DataWrapper[T]) Unpack() (T, error) {
	return dw.Data, dw.Err
}

func (dw DataWrapper[T]) IsOk() bool {
	return dw.Err == nil
}
