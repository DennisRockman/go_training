package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrRuntimeTaskError = errors.New("runtime error in task")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	var errorsChan chan<- error

	if n <= 0 {
		return nil
	}

	if m > 0 {
		errorsChan = make(chan<- error, m)
	}

	tasksChan := make(chan Task, len(tasks))
	for _, task := range tasks {
		tasksChan <- task
	}
	close(tasksChan)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(&wg, tasksChan, errorsChan, m)
	}
	wg.Wait()

	if m > 0 && len(errorsChan) == m {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(wg *sync.WaitGroup, tasksChan <-chan Task, errorsChan chan<- error, errLimit int) {
	defer func() {
		if err := recover(); err != nil && errLimit > 0 {
			writeError(errorsChan, ErrRuntimeTaskError)
		}
		wg.Done()
	}()

	for task := range tasksChan {
		err := task()
		if errLimit > 0 && err != nil && writeError(errorsChan, err) {
			return
		}
	}
}

func writeError(errorsChan chan<- error, err error) bool {
	isReturn := false
	select {
	case errorsChan <- err:
	default:
		isReturn = true
	}
	return isReturn
}
