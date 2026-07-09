package job_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/cubny/lite-reader/internal/infra/job"
)

const tick = 50 * time.Millisecond

// countingJob records how many times it ran, and optionally blocks for a fixed
// duration to simulate a long-running refresh.
type countingJob struct {
	runs  atomic.Int32
	block time.Duration
}

func (c *countingJob) Execute() {
	c.runs.Add(1)
	if c.block > 0 {
		time.Sleep(c.block)
	}
}

func TestScheduler_RunsJobsAtStart(t *testing.T) {
	j := &countingJob{}
	s := job.NewScheduler(time.Hour, j)
	s.Start()
	defer s.Stop()

	assert.Eventually(t, func() bool {
		return j.runs.Load() == 1
	}, time.Second, 5*time.Millisecond, "job should run once at boot without waiting for a tick")
}

// Regression test for the scheduler that stopped firing after the first tick,
// leaving feeds stale indefinitely.
func TestScheduler_RunsJobsOnEveryTick(t *testing.T) {
	j := &countingJob{}
	s := job.NewScheduler(tick, j)
	s.Start()
	defer s.Stop()

	assert.Eventually(t, func() bool {
		return j.runs.Load() >= 4
	}, 2*time.Second, tick/2, "job should keep running on every tick, got %d runs", j.runs.Load())
}

func TestScheduler_RunsAllJobs(t *testing.T) {
	first, second := &countingJob{}, &countingJob{}
	s := job.NewScheduler(tick, first, second)
	s.Start()
	defer s.Stop()

	assert.Eventually(t, func() bool {
		return first.runs.Load() >= 2 && second.runs.Load() >= 2
	}, 2*time.Second, tick/2, "every registered job should run on each tick")
}

// Regression test for the shutdown hang: once a tick had fired, the worker
// parked on a channel it never drained and stopped observing quit, so Stop
// blocked forever and fly.io had to SIGKILL the machine on every deploy.
func TestScheduler_StopReturnsAfterTicksHaveFired(t *testing.T) {
	j := &countingJob{}
	s := job.NewScheduler(tick, j)
	s.Start()

	assert.Eventually(t, func() bool {
		return j.runs.Load() >= 2
	}, 2*time.Second, tick/2, "precondition: at least one tick must fire before Stop")

	stopped := make(chan struct{})
	go func() {
		s.Stop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-time.After(2 * time.Second):
		t.Fatal("Stop deadlocked after a tick had fired")
	}
}

func TestScheduler_StopWaitsForInFlightJob(t *testing.T) {
	j := &countingJob{block: 200 * time.Millisecond}
	s := job.NewScheduler(time.Hour, j)
	s.Start()

	// The boot run is in flight by now; Stop must not cut it short.
	start := time.Now()
	s.Stop()

	assert.GreaterOrEqual(t, time.Since(start), 100*time.Millisecond,
		"Stop should wait for the in-flight job to finish")
	assert.Equal(t, int32(1), j.runs.Load())
}

func TestScheduler_StopIsSafeBeforeAnyTick(t *testing.T) {
	j := &countingJob{}
	s := job.NewScheduler(time.Hour, j)
	s.Start()

	stopped := make(chan struct{})
	go func() {
		s.Stop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-time.After(2 * time.Second):
		t.Fatal("Stop deadlocked before the first tick")
	}
}
