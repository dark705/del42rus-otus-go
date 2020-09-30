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
	testCh := make(chan string, 1000)
	tasks := getTasks(1000, "rnd", testCh)
	_ = Run(tasks, 100, 10)

	n := runtime.NumGoroutine() - 2 //-1 main + -1 test itself
	if n != 0 {
		t.Error("Not all goroutines finished after Run() done by M errors:", n)
	}
}

func TestRunAllGoroutinesFinishedByMErrors(t *testing.T) {
	testCh := make(chan string, 1000)
	tasks := getTasks(1000, "rnd", testCh)
	_ = Run(tasks, 100, 10)

	n := runtime.NumGoroutine() - 2 //-1 main + -1 test itself
	if n != 0 {
		t.Error("Not all goroutines finished after Run() done by M errors:", n)
	}
}

func TestRunAllSuccessDone(t *testing.T) {
	nTask := 10000
	N := 1000
	M := 0

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
