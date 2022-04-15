package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	tests := []struct {
		name           string
		workersCount   int
		maxErrorsCount int
	}{
		{"if were errors in first M tasks, than finished not more N+M tasks", 10, 23},
		{"if m = 0 then ignore errors", 10, 0},
		{"if m < 0 then ignore errors", 10, -1},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			tasksCount := 50
			tasks := make([]Task, 0, tasksCount)

			var runTasksCount int32

			for i := 0; i < tasksCount; i++ {
				err := fmt.Errorf("error from task %d", i)
				tasks = append(tasks, func() error {
					taskSleep := time.Millisecond*time.Duration(rand.Intn(100)) + time.Millisecond
					require.Eventually(t, func() bool { return true }, taskSleep+time.Second, taskSleep)
					atomic.AddInt32(&runTasksCount, 1)
					return err
				})
			}

			workersCount := tst.workersCount
			maxErrorsCount := tst.maxErrorsCount
			err := Run(tasks, workersCount, maxErrorsCount)
			if tst.maxErrorsCount <= 0 {
				require.NoError(t, err)
				require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
			} else {
				require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
				require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
			}
		})
	}

	tests2 := []struct {
		name           string
		workersCount   int
		maxErrorsCount int
		isPanic        bool
	}{
		{"tasks without errors", 5, 1, false},
		{"tasks with panic errors. all tasks were executed", 10, 5, true},
	}
	for _, tst := range tests2 {
		t.Run(tst.name, func(t *testing.T) {
			tasksCount := 50
			tasks := make([]Task, 0, tasksCount)

			var runTasksCount int32
			var sumTime time.Duration

			beginNumber := 0
			taskSleep := time.Millisecond*time.Duration(rand.Intn(100)) + time.Millisecond
			sumTime += taskSleep
			if tst.isPanic {
				beginNumber = 1
				tasks = append(tasks, func() error {
					require.Eventually(t, func() bool { return true }, taskSleep+time.Second, taskSleep)
					atomic.AddInt32(&runTasksCount, 1)
					var divByZero int
					println(5 / divByZero) // one task with panic error
					return nil
				})
			}

			for i := beginNumber; i < tasksCount; i++ {
				taskSleep = time.Millisecond*time.Duration(rand.Intn(100)) + time.Millisecond
				sumTime += taskSleep
				tasks = append(tasks, func() error {
					require.Eventually(t, func() bool { return true }, taskSleep+time.Second, taskSleep)
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}

			workersCount := tst.workersCount
			maxErrorsCount := tst.maxErrorsCount

			start := time.Now()
			err := Run(tasks, workersCount, maxErrorsCount)
			elapsedTime := time.Since(start)
			require.NoError(t, err)
			require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
			require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
		})
	}
}
