package dag

import (
	"context"
	"go-dag/basic"
)

// INode 基础节点接口
type INode[C any] interface {
	GetName() string // 节点名称，唯一性
	GetDependOns() []string
	Get(nCtx C) DataWrapper[any]
}

// DagContext 默认的调度上下文
type DagContext struct {
	Ctx context.Context

	*basic.LogTrace
}
