package main

import (
	"time"

	"go-dag/basic"
	"go-dag/dag"
)

type Task1[C any] struct{}

var task1Instance = &Task1[*dag.DagContext]{}

func Task1Instance() *Task1[*dag.DagContext] {
	return task1Instance
}

func (n *Task1[C]) Get(dCtx *dag.DagContext) dag.DataWrapper[any] {
	time.Sleep(time.Second)
	return dag.DWOf[string]("done")
}

func (n *Task1[C]) GetName() string {
	return basic.GetStructName(n)
}

func (n *Task1[C]) GetDependOns() []string {
	return []string{
		Task2Instance().GetName(), // 依赖 Task1
	}
}
