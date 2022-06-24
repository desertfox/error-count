package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"desertfox.dev/error-count/v1/pkg/category"
	"desertfox.dev/error-count/v1/pkg/worker"
	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/desertfox/gograylog"
	"github.com/go-co-op/gocron"
)

var (
	data       []byte
	err        error
	freq       string = os.Getenv("EC_FREQ")
	webhookUrl string = os.Getenv("EC_TEAMSWEBHOOK")
)

func main() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(freq + "m").Do(do)
	s.StartBlocking()
}

func do() {
	if len(os.Args) > 1 {
		data, err = ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
	} else {
		c := gograylog.New(os.Getenv("EC_HOST"), os.Getenv("EC_USER"), os.Getenv("EC_PASS"))

		f, _ := strconv.Atoi(freq)

		data, err = c.Execute(os.Getenv("EC_QUERY"), os.Getenv("EC_STREAMID"), []string{"message"}, 10000, f)
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

	var results map[string]int = make(map[string]int)
	for r := range wp.Results() {
		if k, ok := results[r.Key]; !ok {
			results[r.Key] = 1

		} else {
			results[r.Key] = k + 1
		}
	}

	keys := make([]string, 0, len(results))
	for key := range results {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return results[keys[i]] > results[keys[j]] })

	keys = keys[0:10]

	var output string
	for k := range keys {
		output = output + fmt.Sprintf("%d:%s\n\r", results[keys[k]], keys[k])
	}

	sendResults(output)

	/*
		var totals map[string]int = storage.Load("./totals.yaml")
		for k, v := range totals {
			fmt.Printf("count:%v\nkey:%v\n\n", v, k)
		}
		err = storage.Save("./totals.yaml", totals)
		if err != nil {
			panic(err)
		}
	*/
}

func sendResults(s string) {
	mstClient := goteamsnotify.NewClient()

	card := goteamsnotify.NewMessageCard()
	card.Title = fmt.Sprintf("Error Counts, Every %smin", freq)
	card.Text = s

	if err := mstClient.Send(webhookUrl, card); err != nil {
		log.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}
}
