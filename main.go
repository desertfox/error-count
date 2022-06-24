package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"desertfox.dev/error-count/v1/pkg/category"
	"desertfox.dev/error-count/v1/pkg/storage"
	"desertfox.dev/error-count/v1/pkg/worker"
	"github.com/desertfox/gograylog"
)

var (
	data []byte
	err  error
)

func main() {
	if len(os.Args) > 1 {
		data, err = ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
	} else {
		c := gograylog.New(os.Getenv("EC_HOST"), os.Getenv("EC_USER"), os.Getenv("EC_PASS"))

		data, err = c.Execute(os.Getenv("EC_QUERY"), os.Getenv("EC_STREAMID"), []string{"message"}, 10000, 60*5)
		if err != nil {
			panic(err)
		}
	}

	lines := strings.Split(string(data), "\n")
	var jobs []worker.Job = make([]worker.Job, len(lines))

	for i := range lines {
		jobs[i] = worker.Job{
			Meta: worker.JobMeta{
				Id:   i,
				Data: []byte(lines[i]),
			},
			ExecFn: category.FileLineKeyFn(),
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
		fmt.Printf("count:%v\nkey:%v\n\n", v, k)
	}

	err = storage.Save("./totals.yaml", totals)
	if err != nil {
		panic(err)
	}

}
