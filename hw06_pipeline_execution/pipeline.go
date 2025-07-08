package hw06pipelineexecution

import (
	"sort"
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

type Result struct {
	value interface{}
	index int
}

func ExecuteStages(in In, done In, sortCh chan<- Result, stages ...Stage) {
	wg := sync.WaitGroup{}
	var i int

	for value := range in {
		wg.Add(1)
		go func(i int) {
			currentValue := value

			for _, stage := range stages {
				ch := make(chan interface{})
				stageCh := stage(ch)

				go func() {
					ch <- currentValue
				}()

				currentValue = <-stageCh
			}

			sortCh <- Result{value: currentValue, index: i}
			wg.Done()
		}(i)
		i++
	}

	wg.Wait()
	close(sortCh)
}

func SortResults(sortCh <-chan Result, resultCh chan<- interface{}) {
	results := []Result{}

	for result := range sortCh {
		results = append(results, result)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	for _, res := range results {
		resultCh <- res.value
	}

	close(resultCh)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	resultCh := make(chan interface{})
	sortCh := make(chan Result)

	//TODO: добавить реагирование на done
	// Написать общую функцию слушания канала select и переиспользовать на всех этапах чтения из каналов
	go ExecuteStages(in, done, sortCh, stages...)
	go SortResults(sortCh, resultCh)

	return resultCh
}
