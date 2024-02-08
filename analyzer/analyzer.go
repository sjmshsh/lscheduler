package analyzer

import "context"

// Result 关系分析结果
type Result struct {
	Relation map[interface{}]map[interface{}][]string
}

// Analyzer 关系分析器
type Analyzer interface {
	AddNode(nodes ...Node)
	Do(ctx context.Context) error
	GetResult() *Result
}
