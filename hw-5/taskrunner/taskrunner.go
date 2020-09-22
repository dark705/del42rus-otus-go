package taskrunner

import (
	"fmt"
)

func Run(tasks []func() error, N int, M int) error {
	L := len(tasks)

	if N > L {
		N = L
	}

	ch := make(chan func() error)
	abort := make(chan struct{})

	go func() {
		defer close(ch)

		select {
		case <-abort:
			return
		default:
			for _, task := range tasks {
				ch <- task
			}
		}
	}()

	errors := make(chan error, L)

	for i := 0; i < N; i++ {
		go func() {
			select {
			case <-abort:
				return
			default:
				for task := range ch {
					errors <- task()
				}
			}
		}()
	}

	var errNum int

	for i := 0; i < L; i++ {
		err := <-errors

		if err != nil {
			errNum++
		}

		if errNum != 0 && errNum == M {
			close(abort)

			return fmt.Errorf("%v tasks returned an error", errNum)
		}
	}

	return nil
}
