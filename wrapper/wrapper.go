package wrapper

import "github.com/sjmshsh/lscheduler/step"

type Wrapper interface {
	step.Step
	GetStep() step.Step
}
