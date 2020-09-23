package taskrunner

import (
	"fmt"
)

func Run(tasks []func() error, N int, M int) error {
	L := len(tasks)

	if N > L {
		N = L
	}

	abort := make(chan struct{})
	aborted := func() bool {
		select {
		case <-abort:
			return true
		default:
			return false
		}
	}

	ch := make(chan func() error)
	var waiters []chan struct{}

	wait := func() {
		for _, waiter := range waiters {
			<-waiter
		}
		<-ch
	}

	go func() {
		defer close(ch)
		for _, task := range tasks {
			if aborted() {
				return
			}
			ch <- task
		}
	}()

	errors := make(chan error, L)
	for i := 0; i < N; i++ {
		waiter := make(chan struct{})
		waiters = append(waiters, waiter)

		go func(waiter chan struct{}) {
			defer close(waiter)
			for task := range ch {
				errors <- task()
			}
		}(waiter)
	}

	var errNum int

	for i := 0; i < L; i++ {
		err := <-errors

		if err != nil {
			errNum++
		}

		if errNum != 0 && errNum >= M {
			close(abort)
			wait()
			return fmt.Errorf("error")
		}
	}

	wait()

	return nil
}
