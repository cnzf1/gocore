package thread_test

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/cnzf1/gocore/thread"
	"github.com/stretchr/testify/assert"
)

func TestRoutinePool(t *testing.T) {
	times := 100
	pool := thread.NewTaskRunner(runtime.NumCPU())

	var counter int32
	var waitGroup sync.WaitGroup
	for i := 0; i < times; i++ {
		waitGroup.Add(1)
		pool.Schedule(func() {
			atomic.AddInt32(&counter, 1)
			waitGroup.Done()
		})
	}

	waitGroup.Wait()

	assert.Equal(t, times, int(counter))
}

func BenchmarkRoutinePool(b *testing.B) {
	r := thread.NewTaskRunner(runtime.NumCPU())
	for i := 0; i < b.N; i++ {
		r.Schedule(func() {
		})
	}
}
