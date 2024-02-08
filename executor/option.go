package executor

type option struct {
	breakOnError bool
}

type Option func(opt *option)

// BreakOnError 遇到错误的时候中断执行, 目前没有实现
func BreakOnError(opt *option) {
	opt.breakOnError = true
}
