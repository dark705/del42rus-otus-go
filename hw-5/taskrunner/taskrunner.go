package taskrunner

import "sync"

func Run(tasks []func() error, N int, M int) {
	if N > len(tasks) {
		N = len(tasks)
	}

	ch := make(chan func() error)

	var wg sync.WaitGroup

	go func() {
		for _, task := range tasks {
			ch <- task
		}

		close(ch)
	}()

	errors := make(chan error, len(tasks))

	for i := 0; i < N; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for task := range ch {
				errors <- task()
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
			return
		}
	}
}
