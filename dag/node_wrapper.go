package dag

// NodeWrapper 节点包装器
type NodeWrapper[C any] struct {
	Node           INode[C]
	IsStrongDepend bool
}

// NodeGroup 节点配置组
type NodeGroup[C any] struct {
	Nws []NodeWrapper[C]
}

func NewNodeGroup[C any](cap int) *NodeGroup[C] {
	return &NodeGroup[C]{Nws: make([]NodeWrapper[C], 0, cap)}
}

func (g *NodeGroup[C]) Append(node INode[C], isStrongDepend bool) {
	g.Nws = append(g.Nws, NodeWrapper[C]{Node: node, IsStrongDepend: isStrongDepend})
}

// DefaultCalcFunc 默认的计算函数
func DefaultCalcFunc[C any](nCtx C, wrapper NodeWrapper[C]) error {
	dw := wrapper.Node.Get(nCtx)
	if dw.Err != nil {
		return dw.Err
	}
	return nil
}
