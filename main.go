package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"desertfox.dev/error-count/v1/pkg/category"
	"desertfox.dev/error-count/v1/pkg/storage"
	"desertfox.dev/error-count/v1/pkg/worker"
	"github.com/desertfox/gograylog"
)

func main() {
	c := gograylog.New(os.Args[1], os.Args[2], os.Args[3])

	data, err := c.Execute("kubernetes_namespace_name:portal* AND error", os.Args[4], []string{"message"}, 10000, 60*15)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")
	var jobs []worker.Job = make([]worker.Job, len(lines))

	for i := range lines {
		jobs[i] = worker.Job{
			Meta: worker.JobMeta{
				Id:   i,
				Type: category.AlphaNumColon,
				Data: []byte(lines[i]),
			},
			ExecFn: category.CreateKeyFn(),
		}
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	wp := worker.NewWP(5)
	go wp.Queue(jobs)
	go wp.Run(ctx)

	var totals map[string]int = storage.Load("./totals.yaml")

	for r := range wp.Results() {
		if k, ok := totals[r.Key]; !ok {
			totals[r.Key] = 1

		} else {
			totals[r.Key] = k + 1
		}
	}

	for k, v := range totals {
		fmt.Printf("category:%v, count:%v\n", k, v)
	}

	err = storage.Save("./totals.yaml", totals)
	if err != nil {
		panic(err)
	}

}
