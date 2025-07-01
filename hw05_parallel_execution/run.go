package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workerCount, errorLimit int) error {
	if workerCount < 1 {
		return nil
	}

	if errorLimit < 1 {
		errorLimit = 1
	}

	taskCh := make(chan Task)
	var errorCount int32

	wg := sync.WaitGroup{}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if atomic.LoadInt32(&errorCount) >= int32(errorLimit) {
					return
				}

				taskErr := task()

				if taskErr != nil {
					atomic.AddInt32(&errorCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errorCount) >= int32(errorLimit) {
			break
		}

		taskCh <- task
	}

	close(taskCh)
	wg.Wait()

	if errorCount >= int32(errorLimit) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
