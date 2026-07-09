package job

import "time"

type Job interface {
	Execute()
}

type Scheduler struct {
	jobs     []Job
	interval time.Duration
	quit     chan struct{}
	done     chan struct{}
}

func NewScheduler(interval time.Duration, jobs ...Job) *Scheduler {
	return &Scheduler{
		jobs:     jobs,
		interval: interval,
		quit:     make(chan struct{}),
		done:     make(chan struct{}),
	}
}

// Start runs every job once, then again on each tick, until Stop is called.
func (s *Scheduler) Start() {
	go func() {
		defer close(s.done)
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		s.runAll()

		for {
			select {
			case <-s.quit:
				return
			case <-ticker.C:
				s.runAll()
			}
		}
	}()
}

func (s *Scheduler) runAll() {
	for _, j := range s.jobs {
		j.Execute()
	}
}

// Stop signals the worker to exit and waits for any in-flight job to finish.
// Safe to call once; subsequent calls will panic on the close of quit.
func (s *Scheduler) Stop() {
	close(s.quit)
	<-s.done
}
