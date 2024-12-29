package errchan

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestOk(t *testing.T) {
	ctx := context.Background()
	ech := WithContext[int](ctx, 0)

	ech.Do(func(ctx context.Context, ch chan<- int) error {
		for i := 0; i <= 3; i++ {
			ch <- i
		}
		return nil
	})

	ech.Do(func(ctx context.Context, ch chan<- int) error {
		for i := 4; i <= 8; i++ {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ch <- i
		}
		return nil
	})

	acc := 0
	for x := range ech.ch {
		acc += x
	}

	if ech.Err() != nil {
		t.Fatal("Error was not expected")
	}

	if acc != 36 {
		t.Fatal("Data is not correct")
	}
}

func TestErr(t *testing.T) {
	ctx := context.Background()
	ech := WithContext[int](ctx, 10)

	ech.Do(func(ctx context.Context, ch chan<- int) error {
		for i := 0; i <= 3; i++ {
			ch <- i
			time.Sleep(10 * time.Millisecond)
			if i == 2 {
				return errors.New("failOne")
			}
		}
		return nil
	})

	ech.Do(func(ctx context.Context, ch chan<- int) error {
		for i := 4; i <= 8; i++ {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ch <- i
			time.Sleep(20 * time.Millisecond)
		}
		return nil
	})

	expErr := errors.New("failFirst")
	ech.Do(func(ctx context.Context, ch chan<- int) error {
		return expErr
	})

	acc := 0
	for x := range ech.ch {
		acc += x
	}

	if ech.Err() != expErr {
		t.Fatal("Error is not correct")
	}

	if acc != 3 {
		t.Fatal("Data is not correct")
	}
}

func TestErr2(t *testing.T) {
	ctx := context.Background()
	ech := WithContext[int](ctx, 10)

	ech.Do(func(ctx context.Context, ch chan<- int) error {
		for i := 0; i <= 3; i++ {
			time.Sleep(10 * time.Millisecond)
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ch <- i
			if i == 2 {
				return errors.New("failOne")
			}
		}
		return nil
	})

	ech.Do(func(ctx context.Context, ch chan<- int) error {
		for i := 4; i <= 8; i++ {
			time.Sleep(20 * time.Millisecond)
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ch <- i
		}
		return nil
	})

	expErr := errors.New("failFirst")
	ech.Do(func(ctx context.Context, ch chan<- int) error {
		return expErr
	})

	acc := 0
	for x := range ech.ch {
		acc += x
	}

	if ech.Err() != expErr {
		t.Fatal("Error is not correct")
	}

	if acc != 0 {
		t.Fatal("Data is not correct")
	}
}
