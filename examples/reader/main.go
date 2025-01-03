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
