package errchan

import (
	"context"
	"sync"
)

type errChan[T any] struct {
	ctx     context.Context
	ch      chan T
	cancel  func()
	err     error
	runOnce sync.Once
	errOnce sync.Once
	wg      sync.WaitGroup
}

func (ech *errChan[T]) Err() error {
	ech.wg.Wait()
	return ech.err
}

func (ech *errChan[T]) Chan() <-chan T {
	ech.done()
	return ech.ch
}

func WithContext[T any](ctx context.Context, bufSize int) *errChan[T] {
	cctx, cancel := context.WithCancel(ctx)

	ech := &errChan[T]{
		ctx:    cctx,
		cancel: cancel,
		ch:     make(chan T, bufSize),
	}
	return ech
}

func (ech *errChan[T]) Go(fn func(ctx context.Context, ch chan<- T) error) {
	ech.wg.Add(1)
	ech.done()
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

func (ech *errChan[T]) done() {
	ech.runOnce.Do(func() {
		go func() {
			ech.wg.Wait()
			close(ech.ch)
		}()
	})
}
