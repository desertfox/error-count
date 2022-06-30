package worker

import (
	"context"
	"sync"

	"desertfox.dev/error-count/v1/pkg/count"
)

type Pool struct {
	instances int
	jobs      chan Job
	results   chan count.Record
}

func NewWP(instances int) Pool {
	return Pool{
		instances: instances,
		jobs:      make(chan Job, instances),
		results:   make(chan count.Record, instances),
	}
}

func (wp Pool) Queue(jobs []Job) {
	for i := range jobs {
		wp.jobs <- jobs[i]
	}
	close(wp.jobs)
}

func (wp Pool) Results() <-chan count.Record {
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

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- count.Record) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			results <- job.execute(ctx)
		case <-ctx.Done():
			results <- count.Record{
				Err: ctx.Err(),
			}
			return
		}
	}
}
