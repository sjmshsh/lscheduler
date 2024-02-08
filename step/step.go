package step

import "context"

// Step 基本流程
type Step interface {
	Name() string
	InputPtr() interface{}
	OutPutPtr() interface{}

	Process(ctx context.Context) error
}
