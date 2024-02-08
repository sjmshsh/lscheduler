package executor

import "context"

// Execution 执行指令
type Execution interface {
	Process(ctx context.Context) error
	Name() string
}
