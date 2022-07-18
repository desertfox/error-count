package worker

import (
	"context"
	"fmt"

	"desertfox.dev/error-count/v1/pkg/count"
)

type Job struct {
	Data   string
	ExecFn FnExec
}

type FnExec func(ctx context.Context, s string) (string, int, error)

func (j Job) execute(ctx context.Context) count.Record {
	file, line, err := j.ExecFn(ctx, j.Data)
	if err != nil {
		fmt.Println(err)
		return count.Record{
			Err: err,
		}
	}

	return count.Record{
		File: file,
		Line: line,
	}
}
