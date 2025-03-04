package main

import (
	"context"
	"fmt"
	"time"

	"go-dag/basic"
	"go-dag/dag"
)

// 测试例子
func main() {
	task1 := Task1Instance()
	task2 := Task2Instance()
	task3 := Task3Instance()

	cfGroup := dag.NewNodeGroup[*dag.DagContext](3)
	cfGroup.Append(task1, true)
	cfGroup.Append(task2, false)
	cfGroup.Append(task3, false)

	customCalcFunc := func(dCtx *dag.DagContext, nw dag.NodeWrapper[*dag.DagContext]) error {
		var err error
		start := time.Now()
		dw := nw.Node.Get(dCtx)
		cost := time.Since(start).Milliseconds()

		if !dw.IsOk() {
			if nw.IsStrongDepend {
				dCtx.LogErr("[%vms][强依赖] pack err!", cost, nw.Node.GetName())
				err = dw.Err
			} else {
				dCtx.LogErr("[%vms][弱依赖] pack err!", cost, nw.Node.GetName())
			}
		} else {
			dCtx.LogOk("[%vms][%v] pack ok!", cost, nw.Node.GetName())
		}

		return err
	}

	scheduler := dag.NewScheduler[*dag.DagContext](&dag.DagContext{Ctx: context.Background(), LogTrace: basic.NewLogTrace()}, cfGroup, customCalcFunc)
	if !scheduler.CircularCheck() {
		panic("circular dependency check failed!")
	}
	scheduler.Run()
	scheduler.PrintDagGraph() // 打印DAG调用图
	fmt.Printf("Trace: %v", basic.ToBeautifulJson(nil, scheduler.DagCtx.LogTrace.Trace.ToSlice()))
}
