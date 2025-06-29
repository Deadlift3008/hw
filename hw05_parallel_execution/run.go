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

func (s *SafeWorkerManager) IncreaseErrorCount() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentErrors++

	if s.currentErrors >= s.errorLimit {
		s.exceedErrorLimit = true
		return true
	}

	return false
}

func (s *SafeWorkerManager) takeNextTask() (Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.currentTaskIndex++

	if s.currentTaskIndex > len(s.tasks)-1 {
		return nil, false
	}

	return s.tasks[s.currentTaskIndex], true
}

func (s *SafeWorkerManager) handle(task *Task) {
	err := (*task)()
	if err != nil {
		exceedErrorLimit := s.IncreaseErrorCount()

		if exceedErrorLimit {
			s.wg.Done()
			return
		}
	}

	nextTask, exist := s.takeNextTask()

	if !exist {
		s.wg.Done()
		return
	}

	s.handle(&nextTask)
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workerCount, errorLimit int) error {
	if workerCount < 1 {
		return nil
	}

	normalizedErrorLimit := errorLimit

	if normalizedErrorLimit < 0 {
		normalizedErrorLimit = 0
	}

	workerManager := &SafeWorkerManager{currentTaskIndex: workerCount - 1, tasks: tasks, errorLimit: normalizedErrorLimit}

	for i := 0; i < workerCount; i++ {
		if i+1 > len(tasks) {
			break
		}

		workerManager.wg.Add(1)
		go workerManager.handle(&tasks[i])
	}

	workerManager.wg.Wait()

	if workerManager.exceedErrorLimit {
		return ErrErrorsLimitExceeded
	}

	return nil
}
