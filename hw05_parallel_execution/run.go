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
	errorLimit       int
	currentTaskIndex int
}

func (s *SafeWorkerManager) IncreaseErrorCount() {
	s.mu.Lock()
	s.currentErrors++
	s.mu.Unlock()
}

func (s *SafeWorkerManager) takeNextTaskIndex() int {

}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workerCount, errorLimit int) error {
	workerManager := &SafeWorkerManager{errorLimit: errorLimit, currentTaskIndex: workerCount}

	for i := 0; i < workerCount; i++ {
		if i+1 > len(tasks) {
			break
		}

		workerManager.wg.Add(1)
		go func(i *int) {
			err := tasks[*i]()

			if err == nil {
				workerManager.IncreaseErrorCount()
			}

			//TODO: дополнительно брать задачу горутиной, реализовать через внешний счетчик который увеличивать (или подумать над другим решением)
			//TODO: проверять на кол-во ошибок перед взятием в работу новой
			//TODO: перенести эту логику в методы SafeWorkerManager(?)

			workerManager.wg.Done()
		}(&i)
	}

	workerManager.wg.Wait()

	return nil
}
