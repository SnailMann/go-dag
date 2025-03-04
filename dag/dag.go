package dag

import (
	"fmt"
	"time"

	"github.com/eapache/queue"
)

// Scheduler 用于调度DAG中的任务
type Scheduler[C any] struct {
	DagCtx     C
	NodeGroup  *NodeGroup[C]
	NameToNode map[string]NodeWrapper[C]

	CalcFunc func(dCtx C, conf NodeWrapper[C]) error // 执行的包装函数，提供给使用方自行定制

	BeReliedOn  map[string][]NodeWrapper[C]
	InDegree    map[string]int
	ExecResults map[string]*RunResponse[C]

	Done chan *RunResponse[C]
}

type RunResponse[C any] struct {
	ExecStartTime time.Time
	ExecEndTime   time.Time
	NodeWrapper   NodeWrapper[C]
	Error         error
}

func NewScheduler[C any](dCtx C, nodeGroup *NodeGroup[C], calcFunc func(dCtx C, conf NodeWrapper[C]) error) *Scheduler[C] {
	res := &Scheduler[C]{
		DagCtx:      dCtx,
		CalcFunc:    calcFunc,
		NodeGroup:   nodeGroup,
		NameToNode:  make(map[string]NodeWrapper[C]),
		InDegree:    make(map[string]int),
		BeReliedOn:  make(map[string][]NodeWrapper[C]),
		ExecResults: make(map[string]*RunResponse[C]),
		Done:        make(chan *RunResponse[C], len(nodeGroup.Nws)),
	}

	for _, c := range nodeGroup.Nws {
		res.NameToNode[c.Node.GetName()] = c
	}

	// 记录节点被谁依赖了
	for _, nw := range res.NodeGroup.Nws {
		for _, beReliedNode := range nw.Node.GetDependOns() {
			if _, exist := res.NameToNode[beReliedNode]; !exist {
				panic(fmt.Sprintf("node %v be relied on not exist", beReliedNode))
			}
			res.BeReliedOn[beReliedNode] = append(res.BeReliedOn[beReliedNode], nw)
		}
	}

	// 记录节点的入度数
	for _, nw := range nodeGroup.Nws {
		res.InDegree[nw.Node.GetName()] = len(nw.Node.GetDependOns())
	}
	return res
}

// CircularCheck 循环依赖检查
func (s *Scheduler[C]) CircularCheck() bool {
	que := queue.New()
	inDegree := s.InDegree

	for _, nw := range s.NodeGroup.Nws {
		if inDegree[nw.Node.GetName()] == 0 {
			que.Add(nw)
		}
	}

	var count int
	for que.Length() > 0 {
		nw, ok := que.Peek().(NodeWrapper[C])
		if !ok {
			return false
		}
		que.Remove()
		count++

		for _, beReliedNode := range s.BeReliedOn[nw.Node.GetName()] {
			inDegree[beReliedNode.Node.GetName()]--
			if inDegree[beReliedNode.Node.GetName()] == 0 {
				que.Add(beReliedNode)
			}
		}
	}

	return count == len(s.NodeGroup.Nws)
}

func (s *Scheduler[C]) Run() {
	start := time.Now()
	defer func() {
		fmt.Printf("DAG exec cost: %vms\n", time.Since(start).Milliseconds())
		close(s.Done)
	}()

	// 先把入度为 0 的加入队列执行
	for _, nw := range s.NodeGroup.Nws {
		if s.InDegree[nw.Node.GetName()] == 0 {
			go s.runNode(nw)
		}
	}
	for i := 0; i < len(s.NodeGroup.Nws); i++ {
		result := <-s.Done
		if result.Error != nil {
			fmt.Printf("[%v] run error: %v\n", result.NodeWrapper.Node.GetName(), result.Error)
		}

		// 入度减1，然后把入度为 0 的，加到队列中，并执行
		for _, beReliedNode := range s.BeReliedOn[result.NodeWrapper.Node.GetName()] {
			s.InDegree[beReliedNode.Node.GetName()]--
			if s.InDegree[beReliedNode.Node.GetName()] == 0 {
				s.runNode(beReliedNode)
			}
		}
	}
}

func (s *Scheduler[C]) runNode(nw NodeWrapper[C]) {
	go func() {
		start := time.Now()
		err := s.CalcFunc(s.DagCtx, nw)
		s.Done <- &RunResponse[C]{
			ExecStartTime: start,
			ExecEndTime:   time.Now(),
			NodeWrapper:   nw,
			Error:         err,
		}
	}()
}
