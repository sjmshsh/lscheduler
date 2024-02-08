package step

import (
	"context"
	"sync"
)

type hookFunc func(ctx context.Context)

// ListenterStep 用于在外部订阅部分字段
type ListenterStep struct {
	*SimpleStep

	listenMutex       sync.RWMutex
	beforeNoticeHooks []hookFunc
	afterNoticeHooks  []hookFunc
}

// NewListenerStep .
func NewListenerStep(name string, concerns interface{}) *ListenterStep {
	var l ListenterStep
	l.SimpleStep = NewSimpleStep(name, concerns, nil)
	l.listenMutex.Lock()
	return &l
}

// Process .
func (l *ListenterStep) Process(ctx context.Context) error {
	for _, h := range l.beforeNoticeHooks {
		h(ctx)
	}

	l.listenMutex.Unlock()

	for _, h := range l.afterNoticeHooks {
		h(ctx)
	}

	return nil
}

// Listen 监听字段，concerns中的字段准备就绪前阻塞
func (l *ListenterStep) Listen() {
	l.listenMutex.RLock() // 无需释放
}

// AddBeforeNoticeHook 添加回调（在Listen方法返回前）
func (l *ListenterStep) AddBeforeNoticeHook(hooks ...hookFunc) {
	l.beforeNoticeHooks = append(l.beforeNoticeHooks, hooks...)
}

// AddAfterNoticeHook 添加回调（在Listen方法返回后）
func (l *ListenterStep) AddAfterNoticeHook(hooks ...hookFunc) {
	l.afterNoticeHooks = append(l.afterNoticeHooks, hooks...)
}
