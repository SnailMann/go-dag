package dag

import (
	"fmt"
	"net/url"
	"strings"
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
		s.ExecResults[result.NodeWrapper.Node.GetName()] = result
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

func (s *Scheduler[C]) PrintDagGraph() {
	var dotBuilder strings.Builder
	dotBuilder.WriteString("digraph G {\n")
	dotBuilder.WriteString("  rankdir=LR;\n")
	dotBuilder.WriteString("  nw [shape=box style=filled];\n\n")

	// 生成节点定义
	for _, nw := range s.NodeGroup.Nws {
		result, exists := s.ExecResults[nw.Node.GetName()]
		if !exists {
			dotBuilder.WriteString(fmt.Sprintf("  %q [label=\"%s\\n(Not Executed)\" fillcolor=gray90];\n",
				nw.Node.GetName(), nw.Node.GetName()))
			continue
		}

		status := "Ok"
		fillColor := "#90EE90" // 浅绿色
		if result.Error != nil {
			status = "Error"
			fillColor = "#FFB6C1" // 浅红色
		}
		duration := result.ExecEndTime.Sub(result.ExecStartTime).Round(time.Millisecond)

		dotBuilder.WriteString(fmt.Sprintf("  %q [label=\"%s\\n%s\\n%s\" fillcolor=\"%s\"];\n",
			nw.Node.GetName(), nw.Node.GetName(), status, duration, fillColor))
	}

	// 生成边定义
	dotBuilder.WriteString("\n")
	for _, nw := range s.NodeGroup.Nws {
		for _, dep := range nw.Node.GetDependOns() {
			depResult, depExists := s.ExecResults[dep]
			if !depExists {
				dotBuilder.WriteString(fmt.Sprintf("  %q -> %q [label=\"(dep not executed)\" style=dashed];\n",
					dep, nw.Node.GetName()))
				continue
			}

			status := "ok"
			color := "green"
			if depResult.Error != nil {
				status = "error"
				color = "red"
			}
			duration := depResult.ExecEndTime.Sub(depResult.ExecStartTime).Round(time.Millisecond)

			dotBuilder.WriteString(fmt.Sprintf("  %q -> %q [label=\"%s (%s)\" color=%s];\n",
				dep, nw.Node.GetName(), status, duration, color))
		}
	}

	dotBuilder.WriteString("}\n")
	// 修正URL编码逻辑
	dotCode := dotBuilder.String()
	// 使用PathEscape代替QueryEscape，并手动处理特殊字符
	encoded := url.PathEscape(dotCode) // 这里会自动将空格转为%20
	graphURL := fmt.Sprintf("https://dreampuf.github.io/GraphvizOnline/#%s", encoded)
	fmt.Printf("graph: %v\n", graphURL)
}
