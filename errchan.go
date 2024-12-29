package errchan

import (
	"context"
	"sync"
)

type ErrChan[T any] struct {
	ctx     context.Context
	cancel  func()
	wg      *sync.WaitGroup
	ch      chan T
	runOnce sync.Once
	errOnce sync.Once
	err     error
}

func (ech *ErrChan[T]) Err() error {
	return ech.err
}

func WithContext[T any](ctx context.Context, bufSize int) *ErrChan[T] {
	cctx, cancel := context.WithCancel(ctx)

	return &ErrChan[T]{
		ctx:    cctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
		ch:     make(chan T, bufSize),
	}
}

func (ech *ErrChan[T]) Do(fn func(ctx context.Context, ch chan<- T) error) {
	ech.wg.Add(1)
	ech.runOnce.Do(func() {
		go func() {
			ech.wg.Wait()
			close(ech.ch)
		}()
	})

	go func() {
		defer ech.wg.Done()
		if err := fn(ech.ctx, ech.ch); err != nil {
			ech.errOnce.Do(func() {
				ech.err = err
				ech.cancel()
			})
		}
	}()
}
