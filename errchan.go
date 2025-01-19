// The package handles work, channels, error-handling and context cancellation under the same hood.
// One of the reasons for the package was to use with buffered channels with slow reader and(or) writer in parallel-async mode.
// It is like the [pkg/golang.org/x/sync/errgroup] with addition channels and context
package errchan

import (
	"context"
	"sync"
)

// The Chan struct creates a channel and a context which can be filled with one or more goroutines working on the same [errchan.Chan]
type Chan[T any] struct {
	ch      chan T
	cancel  func()
	err     error
	runOnce sync.Once
	errOnce sync.Once
	wgGo    sync.WaitGroup
	wgDone  sync.WaitGroup
}

// Creates a new [errchan.Chan] with context.
// bufsize - 0 to make buffered, >0 to make unbuffered channel.
// The context will be cancelled if any related goroutines return an error.
func WithContext[T any](ctx context.Context, bufSize int) (*Chan[T], context.Context) {
	cctx, cancel := context.WithCancel(ctx)

	ech := &Chan[T]{
		cancel: cancel,
		ch:     make(chan T, bufSize),
	}
	return ech, cctx
}

// Starts a goroutine with a function that can return an error and receives a write channel and a context which can be cancelled.
// The first goroutine returns error can be extracted with [Chan.Err] later.
// The channel automatically closes after all related goroutines complete.
func (ech *Chan[T]) Go(fn func(ch chan<- T) error) {
	ech.wgGo.Add(1)
	go func() {
		defer ech.wgGo.Done()
		if err := fn(ech.ch); err != nil {
			ech.errOnce.Do(func() {
				ech.err = err
				ech.cancel()
			})
		}
	}()
}

// Chan returns a reader channel. The channel won't be closed until all the related goroutines have finished
func (ech *Chan[T]) Chan() <-chan T {
	ech.done()
	return ech.ch
}

// Err waits until all goroutines finish and returns an error if any of them return an error, otherwise nil.
func (ech *Chan[T]) Err() error {
	ech.Wait()
	return ech.err
}

// Just waits for related goroutes to finish. [Chan.Err] method calls is too.
// It can be used with defer if you want to be sure that all goroutines are done in case of early return from your function before calling [Chan.Err].
func (ech *Chan[T]) Wait() {
	ech.done()
	ech.wgDone.Wait()
	// for range ech.ch { // TODO: not sure if we need to drain
	// }
}

func (ech *Chan[T]) done() {
	ech.runOnce.Do(func() {
		ech.wgDone.Add(1)
		go func() {
			ech.wgGo.Wait()
			close(ech.ch)
			ech.wgDone.Done()
		}()
	})
}
