package main

import (
	"context"
	"fmt"
	"sync"
)

type Pool struct {
	instances int
	jobs      chan Job
	results   chan Record
}

type Job struct {
	Data   string
	ExecFn FnExec
}

type FnExec func(ctx context.Context, s string) (string, int, error)

func (j Job) execute(ctx context.Context) Record {
	file, line, err := j.ExecFn(ctx, j.Data)
	if err != nil {
		fmt.Println(err)
		return Record{
			Err: err,
		}
	}

	return Record{
		File: file,
		Line: line,
	}
}

func NewWP(instances int) Pool {
	return Pool{
		instances: instances,
		jobs:      make(chan Job, instances),
		results:   make(chan Record, instances),
	}
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
