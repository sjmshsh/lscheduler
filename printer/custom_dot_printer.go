package printer

import (
	"context"
	"fmt"

	"github.com/sjmshsh/lscheduler/wrapper"

	"github.com/awalterschulze/gographviz"
)

// CustomDotNode dot绘制节点
type CustomDotNode interface {
	Name() string
	Attrs() map[string]string
	EdgeAttrs(to DotNode, relation map[interface{}]map[interface{}][]string) map[string]string
}

// CustomDotPrinter 简易dot绘制器
type CustomDotPrinter struct {
}

// NewCustomDotPrinter 新建简易dot绘制器
func NewCustomDotPrinter() Printer {
	return &CustomDotPrinter{}
}

// Print .
func (p CustomDotPrinter) Print(ctx context.Context, relation map[interface{}]map[interface{}][]string) (string, error) {
	graph := gographviz.NewGraph()
	graphAst, _ := gographviz.Parse([]byte(`digraph DAG{rankdir=LR}`))
	_ = gographviz.Analyse(graphAst, graph)

	allEdge := make(map[CustomDotNode]map[CustomDotNode]bool)
	allNode := make(map[CustomDotNode]bool)

	for from, nexts := range relation {
		fromNode, err := p.convertToNode(from)
		if err != nil {
			return "", err
		}

		allNode[fromNode] = true

		if len(nexts) == 0 {
			continue
		}

		for next := range nexts {
			nextNode, err := p.convertToNode(next)
			if err != nil {
				return "", err
			}

			allNode[nextNode] = true

			if len(allEdge[fromNode]) == 0 {
				allEdge[fromNode] = make(map[CustomDotNode]bool)
			}

			allEdge[fromNode][nextNode] = true
		}
	}

	for node := range allNode {
		_ = graph.AddNode("DAG", fmt.Sprintf("\"%s\"", node.Name()), node.Attrs())
	}

	for from, nexts := range allEdge {
		for next := range nexts {
			_ = graph.AddEdge(fmt.Sprintf("\"%s\"", from.Name()), fmt.Sprintf("\"%s\"", next.Name()), true, from.EdgeAttrs(next, relation))
		}
	}

	return graph.String(), nil
}

// Print .
func (p CustomDotPrinter) convertToNode(s interface{}) (CustomDotNode, error) {
	w, ok := s.(wrapper.Wrapper)
	if ok {
		s = w.GetStep()
	}

	n, ok := s.(CustomDotNode)
	if !ok {
		return nil, fmt.Errorf("wrong node type")
	}

	return n, nil
}
