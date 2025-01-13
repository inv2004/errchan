package errchan

import (
	"context"
	"errors"
	"testing"
	"time"
)

func checkErrChan[T any](t *testing.T, ech *Chan[T], expErr error) {
	if ech.Err() != expErr {
		if expErr == nil {
			t.Fatal("Error was not expected")
		} else {
			t.Fatal("Error is not correct")
		}
	}

	_, ok := <-ech.ch
	if ok {
		t.Fatal("Chan was not closed")
	}
}

func TestGoReadErr(t *testing.T) {
	ctx := context.Background()
	ech := WithContext[int](ctx, 0)

	ech.Go(func(ctx context.Context, ch chan<- int) error {
		for i := 0; i <= 3; i++ {
			ch <- i
		}
		return nil
	})

	ech.Go(func(ctx context.Context, ch chan<- int) error {
		for i := 4; i <= 8; i++ {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ch <- i
		}
		return nil
	})

	acc := 0
	for x := range ech.Chan() {
		acc += x
	}

	checkErrChan(t, ech, nil)

	if acc != 36 {
		t.Fatal("Data is not correct")
	}
}

func TestGoReadErr1(t *testing.T) {
	ctx := context.Background()
	ech := WithContext[int](ctx, 0)

	ech.Go(func(ctx context.Context, ch chan<- int) error {
		for i := 0; i <= 3; i++ {
			ch <- i
			time.Sleep(10 * time.Millisecond)
			if i == 2 {
				return errors.New("failOne")
			}
		}
		return nil
	})

	ech.Go(func(ctx context.Context, ch chan<- int) error {
		for i := 4; i <= 8; i++ {
			time.Sleep(50 * time.Millisecond)
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ch <- i
		}
		return nil
	})

	expErr := errors.New("failFirst")
	ech.Go(func(ctx context.Context, ch chan<- int) error {
		return expErr
	})

	acc := 0
	for x := range ech.Chan() {
		acc += x
	}
	checkErrChan(t, ech, expErr)

	if acc != 3 {
		t.Fatal("Data is not correct")
	}
}

func TestGoReadErr2(t *testing.T) {
	ctx := context.Background()
	ech := WithContext[int](ctx, 0)

	ech.Go(func(ctx context.Context, ch chan<- int) error {
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

	ech.Go(func(ctx context.Context, ch chan<- int) error {
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
	ech.Go(func(ctx context.Context, ch chan<- int) error {
		return expErr
	})

	acc := 0
	for x := range ech.Chan() {
		acc += x
	}

	checkErrChan(t, ech, expErr)

	if acc != 0 {
		t.Fatal("Data is not correct")
	}
}

func TestGoWOReadErr(t *testing.T) {
	ctx := context.Background()
	ech := WithContext[int](ctx, 10)

	ech.Go(func(ctx context.Context, ch chan<- int) error {
		for i := 0; i <= 3; i++ {
			time.Sleep(10 * time.Millisecond)
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ch <- i
		}
		return nil
	})

	checkErrChan(t, ech, nil)
}

func TestGoWOReadErr1(t *testing.T) {
	ctx := context.Background()
	ech := WithContext[int](ctx, 10)

	expErr := errors.New("failOne")
	ech.Go(func(ctx context.Context, ch chan<- int) error {
		for i := 0; i <= 3; i++ {
			time.Sleep(10 * time.Millisecond)
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ch <- i
			if i == 2 {
				return expErr
			}
		}
		return nil
	})

	checkErrChan(t, ech, expErr)
}

func TestWOGoReadErr(t *testing.T) {
	ctx := context.Background()
	ech := WithContext[int](ctx, 0)

	acc := 0
	for x := range ech.Chan() {
		acc += x
	}

	checkErrChan(t, ech, nil)

	if acc != 0 {
		t.Fatal("Data is not correct")
	}
}

func TestWOGoWOReadErr(t *testing.T) {
	ctx := context.Background()
	ech := WithContext[int](ctx, 0)

	checkErrChan(t, ech, nil)
}

func TestGoErrRead(t *testing.T) {
	ctx := context.Background()
	ech := WithContext[int](ctx, 10)

	ech.Go(func(ctx context.Context, ch chan<- int) error {
		for i := 0; i <= 3; i++ {
			time.Sleep(10 * time.Millisecond)
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ch <- i
		}
		return nil
	})

	if ech.Err() != nil {
		t.Fatal("Error was not expected")
	}

	acc := 0
	for x := range ech.Chan() {
		t.Fatal("Channel should be closed already")
		acc += x
	}

	if acc != 0 {
		t.Fatalf("Data is not correct: acc = %d", acc)
	}
}

func TestWOGoErrRead(t *testing.T) {
	ctx := context.Background()
	ech := WithContext[int](ctx, 0)

	if ech.Err() != nil {
		t.Fatal("Error was not expected")
	}

	acc := 0
	for x := range ech.Chan() {
		acc += x
	}

	if acc != 0 {
		t.Fatal("Data is not correct")
	}
}

func TestGoDelayReadErr(t *testing.T) {
	ctx := context.Background()
	ech := WithContext[int](ctx, 10)

	ech.Go(func(ctx context.Context, ch chan<- int) error {
		for i := 0; i <= 3; i++ {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ch <- i
		}
		return nil
	})

	time.Sleep(10 * time.Millisecond)
	acc := 0
	for x := range ech.Chan() {
		acc += x
	}

	checkErrChan(t, ech, nil)

	if acc != 6 {
		t.Fatal("Data is not correct")
	}
}

func TestGoDelayReadDelayReadErr(t *testing.T) {
	ctx := context.Background()
	ech := WithContext[int](ctx, 10)

	ech.Go(func(ctx context.Context, ch chan<- int) error {
		for i := 0; i <= 3; i++ {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ch <- i
		}
		return nil
	})

	time.Sleep(10 * time.Millisecond)
	ech.Chan()
	time.Sleep(10 * time.Millisecond)
	acc := 0
	for x := range ech.Chan() {
		acc += x
	}

	checkErrChan(t, ech, nil)

	if acc != 6 {
		t.Fatalf("Data is not correct: %d", acc)
	}
}

func TestGoDelayGo(t *testing.T) {
	ctx := context.Background()
	ech := WithContext[int](ctx, 10)

	ech.Go(func(ctx context.Context, ch chan<- int) error {
		for i := 0; i <= 3; i++ {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ch <- i
		}
		return nil
	})

	time.Sleep(10 * time.Millisecond)

	ech.Go(func(ctx context.Context, ch chan<- int) error {
		for i := 4; i <= 7; i++ {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ch <- i
		}
		return nil
	})

	acc := 0
	for x := range ech.Chan() {
		acc += x
	}

	checkErrChan(t, ech, nil)

	if acc != 28 {
		t.Fatalf("Data is not correct %d", acc)
	}
}

func writer(ctx context.Context) *Chan[int] {
	ech := WithContext[int](ctx, 10)

	ech.Go(func(ctx context.Context, ch chan<- int) error {
		time.Sleep(100 * time.Millisecond)
		for i := 1; i <= 3; i++ {
			ch <- i
		}
		return nil
	})

	return ech
}

func reader(ech *Chan[int]) int {
	acc := 0
	for x := range ech.Chan() {
		if x >= 2 {
			return acc
		}
		acc += x
	}

	return acc
}

func TestGoReturn(t *testing.T) {
	ctx := context.Background()

	ech := writer(ctx)
	defer ech.Close()
	acc := reader(ech)

	// no ech.Err() call here to check if is close chan
	_, ok := <-ech.ch
	if ok {
		t.Fatal("Chan was not closed")
	}
	if acc != 1 {
		t.Fatalf("Data is not correct %d", acc)
	}
}

func TestClose(t *testing.T) {
	ctx := context.Background()

	ech := writer(ctx)
	ech.Close()

	// no ech.Err() call here to check if is close chan
	_, ok := <-ech.ch
	if ok {
		t.Fatal("Chan was not closed")
	}
}
