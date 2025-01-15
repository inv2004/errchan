Make (golang) Channels Great Again

# `errchan.Chan`

The structure covers work with channels, context, cancellations and goroutine error-handling under the same hood

# Motivation
Golang channels are a powerful communication mechanism for goroutines. Unfortunately, there are some challenges in controlling, cancelling and handling errors using channels. 

There are some helpers like https://pkg.go.dev/golang.org/x/sync/errgroup but the `errchan` mod also adds channels and context.

The main motivation is to have create a structure that can be returned as a channel from any function and takes care of all the stuff like creating, closing and error-control for the called goroutine, which can be difficult and tricky if you implement it using just basic types.

# Example:
https://go.dev/play/p/dI40gU9DeZT

```go
package main

import (
    "context"
    "errors"
    "fmt"

    "github.com/inv2004/errchan"
)

func reader(ctx context.Context) *errchan.Chan[int] {
    ech := errchan.WithContext[int](ctx, 10)

    ech.Go(func(ctx context.Context, ch chan<- int) error {
    // ctx is unused here, because just one goroutine
    for i := 1; i <= 3; i++ {
        ch <- i
    }
    return errors.New("readerError")
  })

    return ech
}

func writer(ech *errchan.Chan[int]) (int, error) {
    cnt := 0
    for x := range ech.Chan() {
        cnt++
        if x == 3 {
            return cnt, errors.New("writerError")
        }
    }

    return cnt, nil
}

func main() {
    ctx := context.Background()
    ech := reader(ctx)
    cnt, err := writer(ech)

    fmt.Println(err)       // writerError
    fmt.Println(cnt)       // 3
    fmt.Println(ech.Err()) // readerError
}
```

# TODO
- Add `.SetLimit` to limit number of goroutines
