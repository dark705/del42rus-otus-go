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
	testCh := make(chan string, 100)
	tasks := getTasks(100, "rnd", testCh)
	_ = Run(tasks, 10, 2)

	n := runtime.NumGoroutine() - 2 //-1 main + -1 test itself
	if n != 0 {
		t.Error("Not all goroutines finished after Run() done by M errors:", n)
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
