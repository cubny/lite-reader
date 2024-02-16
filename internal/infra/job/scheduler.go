package job

import "time"

type Job interface {
	Execute()
}

type Scheduler struct {
	Queue    chan Job
	Interval time.Duration
}

func NewScheduler(interval time.Duration) *Scheduler {
	return &Scheduler{
		Queue:    make(chan Job),
		Interval: interval,
	}
}

func (s *Scheduler) Start() {
	go func() {
		ticker := time.NewTicker(s.Interval)

		for {
			select {
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
func (s *Scheduler) ScheduleOnce(duration time.Duration, job Job) {
	go func() {
		time.Sleep(duration)
		s.Queue <- job
	}()
}
