package worker

import "context"

type JobMeta struct {
	Id   int
	Data []byte
}

type FnExec func(ctx context.Context, b []byte) (string, error)

type Job struct {
	Meta   JobMeta
	ExecFn FnExec
}
type Result struct {
	Key string
	Job Job
	Err error
}

func (j Job) execute(ctx context.Context) Result {
	key, err := j.ExecFn(ctx, j.Meta.Data)
	if err != nil {
		return Result{
			Err: err,
			Job: j,
		}
	}

	return Result{
		Key: key,
		Job: j,
	}
}
