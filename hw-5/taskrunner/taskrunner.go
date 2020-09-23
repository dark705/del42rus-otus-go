package taskrunner

import (
	"fmt"
)

func Run(tasks []func() error, N int, M int) error {
	var abort = make(chan struct{})

	aborted := func() bool {
		select {
		case <-abort:
			return true
		default:
			return false
		}
	}

	L := len(tasks)

	if N > L {
		N = L
	}

	ch := make(chan func() error)

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

	var waiters []chan struct{}

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

		if errNum != 0 && errNum == M {
			close(abort)
			for _, waiter := range waiters {
				<-waiter
			}
			<-ch
			return fmt.Errorf("error")
		}
	}

	for _, waiter := range waiters {
		<-waiter
	}
	<-ch

	return nil
}
