package errchan

import (
	"context"
	"iter"
	"sync"
)

type Chan[T any] struct {
	ctx       context.Context
	ch        chan T
	cancel    func()
	err       error
	runOnce   sync.Once
	errOnce   sync.Once
	drainOnce sync.Once
	wgGo      sync.WaitGroup
	wgDone    sync.WaitGroup
}

func (ech *Chan[T]) Close() {
	ech.done()
	ech.drain()
}

func (ech *Chan[T]) Err() error {
	ech.Close()
	return ech.err
}

func (ech *Chan[T]) Chan() iter.Seq[T] {
	return func(yield func(T) bool) {
		ech.done()
		for x := range ech.ch {
			if !yield(x) {
				break
			}
		}
		ech.drain()
	}
}

func WithContext[T any](ctx context.Context, bufSize int) *Chan[T] {
	cctx, cancel := context.WithCancel(ctx)

	ech := &Chan[T]{
		ctx:    cctx,
		cancel: cancel,
		ch:     make(chan T, bufSize),
	}
	return ech
}

func (ech *Chan[T]) Go(fn func(ctx context.Context, ch chan<- T) error) {
	ech.wgGo.Add(1)
	go func() {
		defer ech.wgGo.Done()
		if err := fn(ech.ctx, ech.ch); err != nil {
			ech.errOnce.Do(func() {
				ech.err = err
				ech.cancel()
			})
		}
	}()
}

func (ech *Chan[T]) done() {
	ech.runOnce.Do(func() {
		ech.wgDone.Add(1)
		go func() {
			defer ech.wgDone.Done()
			ech.wgGo.Wait()
			close(ech.ch)
		}()
	})
}

func (ech *Chan[T]) drain() {
	ech.drainOnce.Do(func() {
		ech.wgDone.Wait()
		for range ech.ch { // TODO: not sure if we need to drain
		}
	})
}
