package executor

import "context"

// Executor 执行器
type Executor interface {
	AddExecution(exe Execution, nexts ...Execution)
	Run(ctx context.Context) error
	PrintDeps() string

	hasCircle() error
}
