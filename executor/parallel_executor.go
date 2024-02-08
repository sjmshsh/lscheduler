package executor

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

// ParallelExecutor 并发执行器
type ParallelExecutor struct {
	exe2Nexts map[Execution]map[Execution]bool
	opt       option
}

// NewParallelExecutor 新建并发执行器
func NewParallelExecutor(opts ...Option) Executor {
	e := &ParallelExecutor{
		exe2Nexts: make(map[Execution]map[Execution]bool),
	}

	for _, op := range opts {
		op(&e.opt)
	}

	return e
}

// AddExecution 添加指令
func (e *ParallelExecutor) AddExecution(exe Execution, nexts ...Execution) {
	if len(nexts) == 0 {
		e.exe2Nexts[exe] = nil
		return
	}

	for _, next := range nexts {
		if len(e.exe2Nexts[exe]) == 0 {
			e.exe2Nexts[exe] = make(map[Execution]bool)
		}

		e.exe2Nexts[exe][next] = true
	}
}

// Run 运行执行器
func (e ParallelExecutor) Run(ctx context.Context) error {
	err := e.hasCircle()
	if err != nil {
		return err
	}

	inDegreeM := e.buildInDegree()
	var mutex sync.Mutex
	var startExes []Execution

	for exe, inDgree := range inDegreeM {
		if inDgree == 0 {
			startExes = append(startExes, exe)
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, node := range startExes {
		node := node
		eg.Go(func() error {
			return e.runExecution(ctx, node, inDegreeM, &mutex, eg)
		})
	}

	return eg.Wait()
}

// PrintDeps 格式化指令执行依赖
func (e ParallelExecutor) PrintDeps() string {
	var allRelation []string
	for exe, nexts := range e.exe2Nexts {
		if len(nexts) == 0 {
			allRelation = append(allRelation, fmt.Sprintf("\n%s -> nil ", exe.Name()))
			continue
		}

		for next := range nexts {
			allRelation = append(allRelation, fmt.Sprintf("\n%s -> %s", exe.Name(), next.Name()))
		}
	}

	var res string
	for _, ralation := range allRelation {
		res += ralation
	}

	return res
}

func (e ParallelExecutor) hasCircle() error {
	inDegreeM := e.buildInDegree()

	shouldBreak := false
	for !shouldBreak {
		shouldBreak = true

		for exe, inDegree := range inDegreeM {
			if inDegree != 0 {
				continue
			}

			inDegreeM[exe]--
			shouldBreak = false
			nexts := e.exe2Nexts[exe]

			if len(nexts) == 0 {
				continue
			}

			for next := range nexts {
				inDegreeM[next]--
				// 校验
				if inDegreeM[next] < 0 {
					return fmt.Errorf("unknown error")
				}
			}
		}
	}

	var circleExecution []string
	for exe, in := range inDegreeM {
		if in > 0 {
			circleExecution = append(circleExecution, exe.Name())
		}
	}

	if len(circleExecution) > 0 {
		return fmt.Errorf("Has circle %s", strings.Join(circleExecution, " "))
	}

	return nil
}

func (e ParallelExecutor) buildInDegree() map[Execution]int {
	inDegreeM := make(map[Execution]int)
	for exe := range e.exe2Nexts {
		inDegreeM[exe] = 0
	}

	for _, nexts := range e.exe2Nexts {
		if len(nexts) == 0 {
			continue
		}

		for next := range nexts {
			inDegreeM[next]++
		}
	}
	return inDegreeM
}

func (e ParallelExecutor) runExecution(ctx context.Context, exe Execution, inDegreeM map[Execution]int, mutex *sync.Mutex, eg *errgroup.Group) error {
	err := exe.Process(ctx)
	if err != nil {
		if e.opt.breakOnError {
			return err
		}
	}

	nexts := e.exe2Nexts[exe]

	mutex.Lock()
	defer mutex.Unlock()

	// 如果该执行器有nexts
	if len(nexts) > 0 {
		for next := range nexts {
			inDegreeM[next]--
			if inDegreeM[next] == 0 {
				next := next
				eg.Go(func() error {
					return e.runExecution(ctx, next, inDegreeM, mutex, eg)
				})
			}
		}
	}

	return nil
}
