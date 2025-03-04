package main

import (
	"time"

	"go-dag/basic"
	"go-dag/dag"
)

type Task2[C any] struct{}

var task2Instance = &Task2[*dag.DagContext]{}

func Task2Instance() *Task2[*dag.DagContext] {
	return task2Instance
}

func (n *Task2[C]) Get(dCtx *dag.DagContext) dag.DataWrapper[any] {
	time.Sleep(time.Second)
	return dag.DWEmpty[int64]()
}

func (n *Task2[C]) GetName() string {
	return basic.GetStructName(n)
}

func (n *Task2[C]) GetDependOns() []string {
	return []string{}
}
