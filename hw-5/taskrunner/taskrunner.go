package taskrunner

import (
	"fmt"
)

var abort = make(chan struct{})

func Run(tasks []func() error, N int, M int) error {
	L := len(tasks)

	if N > L {
		N = L
	}

	ch := make(chan func() error)

	go func() {
		defer func() {
			close(ch)
		}()

		for _, task := range tasks {
			if aborted() {
				return
			}

			ch <- task
		}

		return
	}()

	errors := make(chan error, L)

	var waiters []chan struct{}

	for i := 0; i < N; i++ {
		waiter := make(chan struct{})
		waiters = append(waiters, waiter)

		go func(waiter chan struct{}) {
			defer func() {
				close(waiter)
			}()

			for {
				select {
				case <-abort:
					return
				case task, ok := <-ch:
					if !ok {
						return
					}
					errors <- task()
				}
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
			return fmt.Errorf("error")
		}
	}

	for _, waiter := range waiters {
		<-waiter
	}

	return nil
}

func aborted() bool {
	select {
	case <-abort:
		return true
	default:
		return false
	}
}
