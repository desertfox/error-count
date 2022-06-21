package worker

import (
	"context"
	"sync"
)

type Pool struct {
	instances int
	jobs      chan Job
	results   chan Result
}

func NewWP(instances int) Pool {
	return Pool{
		instances: instances,
		jobs:      make(chan Job, instances),
		results:   make(chan Result, instances),
	}
}

func (wp Pool) Queue(jobs []Job) {
	for i := range jobs {
		wp.jobs <- jobs[i]
	}
	close(wp.jobs)
}

func (wp Pool) Results() <-chan Result {
	return wp.results
}

func (wp Pool) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < wp.instances; i++ {
		wg.Add(1)
		go worker(ctx, &wg, wp.jobs, wp.results)
	}
	wg.Wait()
	close(wp.results)
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			results <- job.execute(ctx)
		case <-ctx.Done():
			results <- Result{
				Err: ctx.Err(),
			}
			return
		}
	}
}
