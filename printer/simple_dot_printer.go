package printer

import (
	"context"
	"fmt"

	"github.com/awalterschulze/gographviz"
)

// DotNode dot绘制节点
type DotNode interface {
	Name() string
}

// SimpleDotPrinter 简易dot绘制器
type SimpleDotPrinter struct {
}

// NewSimpleDotPrinter 新建简易dot绘制器
func NewSimpleDotPrinter() Printer {
	return &SimpleDotPrinter{}
}

// Print .
func (p SimpleDotPrinter) Print(ctx context.Context, relation map[interface{}]map[interface{}][]string) (string, error) {
	graph := gographviz.NewGraph()
	graphAst, _ := gographviz.Parse([]byte(`digraph DAG{rankdir=LR}`))
	_ = gographviz.Analyse(graphAst, graph)

	allEdge := make(map[DotNode]map[DotNode]bool)
	allNode := make(map[DotNode]bool)

	for from, nexts := range relation {
		fromNode := from.(DotNode)
		allNode[fromNode] = true

		if len(nexts) == 0 {
			continue
		}

		for next := range nexts {
			nextNode := next.(DotNode)
			allNode[nextNode] = true

			if len(allEdge[fromNode]) == 0 {
				allEdge[fromNode] = make(map[DotNode]bool)
			}

			allEdge[fromNode][nextNode] = true
		}
	}

	for node := range allNode {
		_ = graph.AddNode("DAG", fmt.Sprintf("\"%s\"", node.Name()), map[string]string{"label": fmt.Sprintf("\"%s\"", node.Name())})
	}

	for from, nexts := range allEdge {
		for next := range nexts {
			_ = graph.AddEdge(fmt.Sprintf("\"%s\"", from.Name()), fmt.Sprintf("\"%s\"", next.Name()), true, map[string]string{})
		}
	}

	return graph.String(), nil
}
