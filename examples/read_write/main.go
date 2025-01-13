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
