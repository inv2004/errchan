Make (golang) Channels Great Again

# `errchan.Chan`

The structure covers work with channels, context, cancellations and goroutine error-handling under the same hood

# Motivation
Golang channels are a powerful communication mechanism for goroutines. Unfortunately, there are some challenges in controlling, cancelling and handling errors using channels. 

There are some helpers like https://pkg.go.dev/golang.org/x/sync/errgroup but the `errchan` mod also adds channels and context.

The main motivation is to have create a structure that can be returned as a channel from any function and takes care of all the stuff like creating, closing and error-control for the called goroutine, which can be difficult and tricky if you implement it using just basic types.

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
