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

	ch := make(chan func() error)

	//Мы за ранее знаем сколько будет таких каналов,
	//по этому лучше сразу задать capacity,
	//иначе waiters = append(waiters, waiter) ведёт к лишним аллокациям памяти
	waiters := make([]chan struct{}, 0, N)

	wait := func() {
		// нет конкурентности, т.е мы по ОЧЕРЕДИ ждём когда отработают все каналы
		// например первый воркер отрабатывал дольше всех, тогда все остальные будут ждать только его
		// Именно по этому валится тест TestRunExecuteNotMoreNPlusM
		for _, waiter := range waiters {
			<-waiter
		}
	}

	go func() {
		defer close(ch)
		for _, task := range tasks {
			select { //Можно проще, например так:
			case <-abort:
				return
			default:
				ch <- task
			}
		}
	}()

	errors := make(chan error, L)
	for i := 0; i < N; i++ {
		waiter := make(chan struct{}) //Вообще логика не совсем верная. Не делают на каждую рутину по каналу
		waiters = append(waiters, waiter)

		go func(waiter chan struct{}) {
			defer close(waiter)
			for task := range ch {
				errors <- task()
			}
		}(waiter)
	}

	var errNum int

	for i := 0; i < L; i++ { //Интересное решение
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
