/*
 * @Author: cnzf1
 * @Date: 2021-08-18 18:07:12
 * @LastEditors: cnzf1
 * @LastEditTime: 2021-11-26 15:29:06
 * @Description: 控制超时时间运行函数
 */
package funx

import (
	"context"
	"time"
)

var (
	// ErrCanceled is the error returned when the context is canceled.
	ErrCanceled = context.Canceled
	// ErrTimeout is the error returned when the context's deadline passes.
	ErrTimeout = context.DeadlineExceeded
)

// DoOption defines the method to customize a DoWithTimeout call.
type DoOption func() context.Context

// DoWithTimeout runs fn with timeout control.
func DoWithTimeout(fn func() error, timeout time.Duration, opts ...DoOption) error {
	parentCtx := context.Background()
	for _, opt := range opts {
		parentCtx = opt()
	}
	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()

	done := make(chan error, 1)
	panicChan := make(chan interface{}, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		done <- fn()
	}()

	select {
	case p := <-panicChan:
		panic(p)
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WithContext customizes a DoWithTimeout call with given ctx.
func WithContext(ctx context.Context) DoOption {
	return func() context.Context {
		return ctx
	}
}
