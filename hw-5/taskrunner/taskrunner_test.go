package taskrunner

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestRunAllGoroutinesFinishedByDoneAllTasks(t *testing.T) {
	testCh := make(chan string, 3)
	tasks := getTasks2(testCh)
	_ = Run(tasks, 3, 1)

	n := runtime.NumGoroutine() - 2 //-1 main + -1 test itself
	if n != 0 {
		t.Error("Not all goroutines finished after Run() done by M errors:", n)
	}
}

func TestRunAllSuccessDone(t *testing.T) {
	nTask := 10000
	N := 100
	M := 10

	testCh := make(chan string, nTask)
	tasks := getTasks(nTask, "suc", testCh)
	res := Run(tasks, N, M)

	if res != nil {
		t.Error("Run() all success tasks return not nil")
	}

	if len(testCh) != nTask {
		t.Error("Not all success task done")
	}
}

func TestRunExecuteNotMoreNPlusM(t *testing.T) {
	nTask := 10000
	N := 100
	M := 10

	testCh := make(chan string, nTask)
	tasks := getTasks(nTask, "err", testCh)
	_ = Run(tasks, N, M)

	if len(testCh) > N+M {
		t.Error("Run() executed more then N + M tasks")
	}
}

func TestRunReturnErrorByMLimit(t *testing.T) {
	nTask := 10
	N := 1000
	M := 10

	testCh := make(chan string, nTask)
	tasks := getTasks(nTask, "err", testCh)
	res := Run(tasks, N, M)

	if res == nil {
		t.Error("Run() error tasks not return error by M limit")
	}
}

func getTasks2(testCh chan string) []func() error {

	tasks := []func() error{func() error {
		time.Sleep(1 * time.Second)
		testCh <- fmt.Sprintf("Task #%d done work\n", 1)

		return fmt.Errorf("error")
	}, func() error {
		time.Sleep(5 * time.Second)
		testCh <- fmt.Sprintf("Task #%d done work\n", 2)

		return nil
	}, func() error {
		time.Sleep(2 * time.Second)

		testCh <- fmt.Sprintf("Task #%d done work\n", 3)

		return nil
	}}

	return tasks
}

func getTasks(num int, typeTask string, testCh chan string) []func() error {
	tasks := make([]func() error, 0, num)
	for i := 0; i < num; i++ {
		tasks = append(tasks, getTask(typeTask, testCh))
	}

	return tasks
}

func getTask(t string, testCh chan string) func() error {
	var task func() error
	switch t {
	case "err":
		task = func() error {
			taskNum := rand.Uint32()
			time.Sleep(time.Millisecond * time.Duration(20+rand.Intn(50)))
			errorMessage := "Error on execute Task#" + strconv.Itoa(int(taskNum))
			testCh <- errorMessage
			return errors.New(errorMessage)
		}
	case "suc":
		task = func() error {
			taskNum := rand.Uint32()
			time.Sleep(time.Millisecond * time.Duration(20+rand.Intn(50)))
			testCh <- fmt.Sprintf("Task #%d done work\n", taskNum)
			return nil
		}
	case "rnd":
		fallthrough
	default:
		task = func() error {
			taskNum := rand.Uint32()
			time.Sleep(time.Millisecond * time.Duration(20+rand.Intn(50)))
			if rand.Intn(100) < 50 {
				errorMessage := "Error on execute Task#" + strconv.Itoa(int(taskNum))
				testCh <- errorMessage
				return errors.New(errorMessage)
			}
			testCh <- fmt.Sprintf("Task #%d done work\n", taskNum)
			return nil
		}
	}

	return task
}
