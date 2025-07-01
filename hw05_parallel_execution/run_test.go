package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("should run tasks concurrently", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		workersCount := 10

		release := make(chan struct{})
		var (
			activeWorkers    int32
			maxActiveWorkers int32
		)

		wg := sync.WaitGroup{}

		for i := 0; i < tasksCount; i++ {
			wg.Add(1)
			tasks = append(tasks, func() error {
				current := atomic.AddInt32(&activeWorkers, 1)

				for {
					oldMax := atomic.LoadInt32(&maxActiveWorkers)

					if current > oldMax {
						if atomic.CompareAndSwapInt32(&maxActiveWorkers, oldMax, current) {
							break
						}
					} else {
						break
					}
				}

				<-release
				atomic.AddInt32(&activeWorkers, -1)
				wg.Done()

				return nil
			})
		}

		var err error

		go func() {
			err = Run(tasks, workersCount, 3)
		}()

		time.Sleep(time.Millisecond * time.Duration(100))

		close(release)
		wg.Wait()

		require.Equal(t, int32(workersCount), maxActiveWorkers)
		require.Nil(t, err, "Got error - %v", err)
	})

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("edge cases", func(t *testing.T) {
		t.Run("tasksCount less than workersCount", func(t *testing.T) {
			tasksCount := 5
			workersCount := 10
			tasks := make([]Task, 0, tasksCount)

			var runTasksCount int32

			for i := 0; i < tasksCount; i++ {
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}

			err := Run(tasks, workersCount, 1)

			require.Nil(t, err, "Got error - %v", err)
			require.Equal(t, runTasksCount, int32(tasksCount))
		})

		t.Run("workersCount less than 1", func(t *testing.T) {
			tasksCount := 2
			workersCount := 0
			tasks := make([]Task, 0, tasksCount)

			var runTasksCount int32

			for i := 0; i < tasksCount; i++ {
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}

			err := Run(tasks, workersCount, 1)

			require.Nil(t, err, "Got error - %v", err)
			require.Equal(t, runTasksCount, int32(0))
		})

		t.Run("errorsLimit less than 1 must handled as 1", func(t *testing.T) {
			tasksCount := 10
			workersCount := 2
			tasks := make([]Task, 0, tasksCount)

			var runTasksCount int32

			for i := 0; i < tasksCount; i++ {
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}

			err := Run(tasks, workersCount, 0)

			require.Nil(t, err, "Got error - %v", err)
			require.Equal(t, runTasksCount, int32(10))
		})
	})
}
