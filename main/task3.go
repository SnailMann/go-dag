package main

import (
	"time"

	"go-dag/basic"
	"go-dag/dag"
)

type Task3[C any] struct{}

var task3Instance = &Task3[*dag.DagContext]{}

func Task3Instance() *Task3[*dag.DagContext] {
	return task3Instance
}

func (n *Task3[C]) Get(dCtx *dag.DagContext) dag.DataWrapper[any] {
	time.Sleep(time.Second)
	return dag.DWOf[bool](true)
}

func (n *Task3[C]) GetName() string {
	return basic.GetStructName(n)
}

func (n *Task3[C]) GetDependOns() []string {
	return []string{}
}
