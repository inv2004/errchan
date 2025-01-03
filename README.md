# errchan
Make (golang) Channels Great Again

Helper structure provides which covers channel, context and goroutine error-handling under the one hood

# Motivation
Golang channels are powerful communication mechanism between goroutines, unfortunately there are few challenges how to control, cancel and handle errors using channels

There some helpers for like https://pkg.go.dev/golang.org/x/sync/errgroup, but the `errchan` mod adds channel and context also.

The main motivation is to have a structure which can be returned like a channel from any function and take care of all the stuff like: create, close, drain and error-control of the called goroutines which can be not easy.

# Example:
```go
import (
	"context"
	"errors"
	"fmt"

	"github.com/inv2004/errchan"
)

func reader(ctx context.Context) *errchan.Chan[int] {
	ech := errchan.WithContext[int](ctx, 10)

	ech.Go(func(ctx context.Context, ch chan<- int) error {
		for i := 0; i <= 3; i++ {
			ch <- i
		}
		return errors.New("fail")
	})

	return ech
}

func main() {
	ctx := context.Background()
	ech := reader(ctx)

	acc := 0
	for x := range ech.Chan() {
		acc += x
	}

	fmt.Println(acc)       // 6
	fmt.Println(ech.Err()) // fail
}
```
