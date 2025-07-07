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

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	resultCh := make(chan interface{})
	orderCh := make(chan Result)
	wg := sync.WaitGroup{}

	//TODO: добавить реагирование на done
	go func() {
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

				orderCh <- Result{value: currentValue, index: i}
				wg.Done()
			}(i)
			i++
		}

		wg.Wait()
		close(orderCh)
	}()

	go func() {
		results := []Result{}

		for result := range orderCh {
			results = append(results, result)
		}

		sort.Slice(results, func(i, j int) bool {
			return results[i].index < results[j].index
		})

		for _, res := range results {
			resultCh <- res.value
		}

		close(resultCh)
	}()

	return resultCh
}
