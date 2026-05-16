package job

import "time"

type Job interface {
	Execute()
}

type Scheduler struct {
	Queue    chan Job
	Interval time.Duration
	quit     chan struct{}
	done     chan struct{}
}

func NewScheduler(interval time.Duration) *Scheduler {
	return &Scheduler{
		Queue:    make(chan Job),
		Interval: interval,
		quit:     make(chan struct{}),
		done:     make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	go func() {
		defer close(s.done)
		ticker := time.NewTicker(s.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-s.quit:
				return
			case job := <-s.Queue:
				job.Execute()
			case <-ticker.C:
				for job := range s.Queue {
					job.Execute()
				}
			}
		}
	}()
}

// Stop signals the worker to exit and waits for any in-flight job to finish.
// Safe to call once; subsequent calls will panic on the close of quit.
func (s *Scheduler) Stop() {
	close(s.quit)
	<-s.done
}
func (s *Scheduler) ScheduleOnce(duration time.Duration, job Job) {
	go func() {
		time.Sleep(duration)
		s.Queue <- job
	}()
}
