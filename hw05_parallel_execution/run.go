package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type SafeWorkerManager struct {
	mu               sync.Mutex
	wg               sync.WaitGroup
	currentErrors    int
	currentTaskIndex int
	tasks            []Task
	errorLimit       int
	exceedErrorLimit bool
}

func (s *SafeWorkerManager) IncreaseErrorCount() {
	s.mu.Lock()
	s.currentErrors++
	s.mu.Unlock()
}

func (s *SafeWorkerManager) takeNextTaskIndex() int {
	s.mu.Lock()
	s.currentTaskIndex++
	s.mu.Unlock()
	return s.currentTaskIndex
}

func (s *SafeWorkerManager) handle(i *int) {
	err := s.tasks[*i]()

	if err != nil {
		s.IncreaseErrorCount()
	}

	if s.currentErrors >= s.errorLimit {
		s.wg.Done()
		s.exceedErrorLimit = true
		return
	}

	nextTaskIndex := s.takeNextTaskIndex()

	if nextTaskIndex > len(s.tasks)-1 {
		s.wg.Done()
		return
	}

	s.handle(&nextTaskIndex)
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workerCount, errorLimit int) error {
	if workerCount < 1 {
		return nil
	}

	workerManager := &SafeWorkerManager{currentTaskIndex: workerCount - 1, tasks: tasks, errorLimit: errorLimit}

	for i := 0; i < workerCount; i++ {
		if i+1 > len(tasks) {
			break
		}

		workerManager.wg.Add(1)
		go workerManager.handle(&i)
	}

	workerManager.wg.Wait()

	if workerManager.exceedErrorLimit {
		return ErrErrorsLimitExceeded
	} else {
		return nil
	}
}
