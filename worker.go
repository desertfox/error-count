package main

import (
	"context"
	"sync"
)

type Pool struct {
	instances int
	jobs      chan Job
	results   chan Record
}

type Job struct {
	Data string
	Fnc  FnExec
}

type FnExec func(ctx context.Context, s string) Record

func (j Job) execute(ctx context.Context) Record {
	return j.Fnc(ctx, j.Data)
}

func (wp Pool) Queue(jobs []Job) {
	for i := range jobs {
		wp.jobs <- jobs[i]
	}
	close(wp.jobs)
}

func (wp Pool) Results() <-chan Record {
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

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Record) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			results <- job.execute(ctx)
		case <-ctx.Done():
			results <- Record{
				Err: ctx.Err(),
			}
			return
		}
	}
}
