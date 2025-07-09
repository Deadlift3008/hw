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

func readValueChannel(channelToRead In, doneChannel In) (interface{}, bool, bool) {
	select {
	case <-doneChannel:
		return nil, false, true
	default:
		select {
		case value, ok := <-channelToRead:
			return value, ok, false
		case <-doneChannel:
			return nil, false, true
		}
	}
}

func readResultValueChannel(channelToRead <-chan Result, doneChannel In) (*Result, bool, bool) {
	select {
	case <-doneChannel:
		return &Result{}, false, true
	default:
		select {
		case value, ok := <-channelToRead:
			return &value, ok, false
		case <-doneChannel:
			return &Result{}, false, true
		}
	}
}

func writeValueChannel(channelToWrite chan interface{}, doneChannel In, valueToWrite interface{}) {
	select {
	case channelToWrite <- valueToWrite:
	case <-doneChannel:
	}
}

func executeStages(in In, done In, sortCh chan<- Result, stages ...Stage) {
	wg := sync.WaitGroup{}
	var i int
	defer wg.Wait()
	defer close(sortCh)

	for {
		value, isOpened, isDone := readValueChannel(in, done)

		if !isOpened || isDone {
			return
		}

		wg.Add(1)
		go func(i int) {
			currentValue := value
			defer wg.Done()

			for _, stage := range stages {
				ch := make(chan interface{})
				stageCh := stage(ch)

				go func() {
					writeValueChannel(ch, done, currentValue)
				}()

				newValue, _, isDone := readValueChannel(stageCh, done)

				if isDone {
					return
				}

				currentValue = newValue
			}

			select {
			case <-done:
				return
			case sortCh <- Result{value: currentValue, index: i}:
			}
		}(i)
		i++
	}
}

func sortResults(sortCh <-chan Result, resultCh chan<- interface{}, done In) {
	results := []Result{}
	defer close(resultCh)

	for {
		result, isOpened, isDone := readResultValueChannel(sortCh, done)

		if !isOpened {
			break
		}

		if isDone {
			return
		}

		results = append(results, *result)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	for _, res := range results {
		select {
		case <-done:
			return
		default:
			select {
			case <-done:
				return
			case resultCh <- res.value:
			}
		}
	}
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	resultCh := make(chan interface{})
	sortCh := make(chan Result)

	go executeStages(in, done, sortCh, stages...)
	go sortResults(sortCh, resultCh, done)

	return resultCh
}
