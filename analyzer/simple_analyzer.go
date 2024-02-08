package analyzer

import (
	"context"
	"fmt"
)

// SimpleAnalyzer 简易分析器
type SimpleAnalyzer struct {
	nodes map[Node]bool
	res   *Result
}

func NewSimpleAnalyzer() Analyzer {
	return &SimpleAnalyzer{
		nodes: make(map[Node]bool),
	}
}

func (a *SimpleAnalyzer) AddNode(nodes ...Node) {
	for _, node := range nodes {
		a.nodes[node] = true
	}
}

func (a *SimpleAnalyzer) Do(ctx context.Context) error {
	output2Nodes := make(map[string]Node)

	for node := range a.nodes {
		// 拿到一个节点的对外输出
		outputs := node.Outputs()
		for _, output := range outputs {
			if _, ok := output2Nodes[output]; ok {
				return fmt.Errorf("conflict output %s", output)
			}

			output2Nodes[output] = node
		}
	}

	relation := make(map[interface{}]map[interface{}][]string)
	for node := range a.nodes {
		// 得到一个节点输入的数据
		inputs := node.Inputs()
		for _, input := range inputs {
			// 某一个节点的输出对应另外一个节点的输入, 得到依赖关系
			dep, ok := output2Nodes[input]
			if !ok {
				return fmt.Errorf("no output for %s", input)
			}

			// 初始化map
			if len(relation[dep]) == 0 {
				relation[dep] = make(map[interface{}][]string)
			}

			relation[dep][node] = append(relation[dep][node], input)
		}
	}

	for node := range a.nodes {
		if len(relation[node]) == 0 {
			relation[node] = nil
		}
	}

	a.res = &Result{Relation: relation}
	return nil
}

func (a *SimpleAnalyzer) GetResult() *Result {
	return a.res
}
