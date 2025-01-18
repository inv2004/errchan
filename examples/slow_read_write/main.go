package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/inv2004/errchan"
)

func slowReader(ctx context.Context) *errchan.Chan[int] {
	ech, ctx := errchan.WithContext[int](ctx, 10)

	ech.Go(func(ch chan<- int) error {
		for i := 1; i <= 3; i++ {
			time.Sleep(1 * time.Second)
			ch <- i
		}
		return errors.New("readerFail")
	})

	return ech
}

func main() {
	ctx := context.Background()
	ech := slowReader(ctx)

	acc := 0
	n := time.Now()
	for x := range ech.Chan() {
		acc += x
		time.Sleep(1 * time.Second) // slow writer
	}
	fmt.Println(time.Since(n))

	fmt.Println(acc)       // 6
	fmt.Println(ech.Err()) // readerFail
}
