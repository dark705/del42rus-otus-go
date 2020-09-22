package taskrunner

import (
	"fmt"
	"sync"
)

func Run(tasks []func() error, N int, M int) error {
	L := len(tasks)

	if N > L {
		N = L
	}

	ch := make(chan func() error)
	abort := make(chan struct{})

	go func() {
		for _, task := range tasks {
			ch <- task
		}

		close(ch)
	}()

	errors := make(chan error, L)

	var wg sync.WaitGroup

	for i := 0; i < N; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

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

	go func() {
		wg.Wait()
		close(errors)
	}()

	var errNum int

	for err := range errors {
		if err != nil {
			errNum++
		}

		if errNum == M {
			close(abort)
			return fmt.Errorf("%v tasks returned an error", errNum)
		}
	}

	return nil
}
