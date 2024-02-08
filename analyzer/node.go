package analyzer

// Node 基本节点
type Node interface {
	Inputs() []string
	Outputs() []string
}
