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

type DoneState struct {
	doneChannel In
	done        bool
	mu          sync.Mutex
}

func (d *DoneState) startListen() {
	go func() {
		<-d.doneChannel
		d.mu.Lock()
		d.done = true
		d.mu.Unlock()
	}()
}

func (d *DoneState) isDone() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.done
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

func executeStages(in In, done In, sortCh chan<- Result, stages ...Stage) {
	wg := &sync.WaitGroup{}
	var i int
	defer close(sortCh)

	doneState := &DoneState{doneChannel: done}
	doneState.startListen()

	for {
		value, isOpened, isDone := readValueChannel(in, done)

		if !isOpened || isDone {
			break
		}

		wg.Add(1)
		go func(i int, v interface{}) {
			currentValue := v
			defer wg.Done()

			for _, stage := range stages {
				if doneState.isDone() {
					return
				}
				ch := make(chan interface{})
				stageCh := stage(ch)

				go func() {
					defer close(ch)

					select {
					case <-done:
						return
					default:
						select {
						case ch <- currentValue:
						case <-done:
							return
						}
					}
				}()

				if doneState.isDone() {
					return
				}

				newValue, ok := <-stageCh

				if !ok {
					continue
				}

				currentValue = newValue
			}

			if doneState.isDone() {
				return
			}

			sortCh <- Result{value: currentValue, index: i}
		}(i, value)
		i++
	}

	wg.Wait()
}

func sortResults(sortCh <-chan Result, resultCh chan<- interface{}, done In) {
	results := []Result{}
	defer close(resultCh)

	for {
		result, isOpened, isDone := readResultValueChannel(sortCh, done)

		if isDone {
			return
		}

		if !isOpened {
			break
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
