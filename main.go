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
	freq       string           = os.Getenv("EC_FREQ")
	webhookUrl string           = os.Getenv("EC_TEAMSWEBHOOK")
	history    []map[string]int = make([]map[string]int, 0)
	totals     map[string]int   = make(map[string]int, 10)
)

func main() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(freq + "m").Do(do)
	s.Every("60m").Do(report)
	s.StartBlocking()
}

func sortKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return m[keys[i]] > m[keys[j]] })

	return keys
}

func report() {
	var hourTotals map[string]int = make(map[string]int)
	for i := range history {
		for file := range history[i] {
			if v, ok := hourTotals[file]; !ok {
				hourTotals[file] = history[i][file]
			} else {
				hourTotals[file] = history[i][file] + v
			}
		}
	}

	sortedKeys := sortKeys(hourTotals)
	var output string = "COUNT_FILE\n\r"
	for k := range sortedKeys {
		output = output + fmt.Sprintf("%03d_%s\n\r", hourTotals[sortedKeys[k]], sortedKeys[k])
	}

	sendResults("Error Count Hour Totals", output)

	history = make([]map[string]int, 0)
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

	sortedKeys := sortKeys(results)
	sortedKeys = sortedKeys[0:10]

	newTotals := make(map[string]int, 10)
	var output string = "COUNT:PREV:+/-:FILE\n\r"
	for k := range sortedKeys {
		var last, change int = 0, results[sortedKeys[k]]
		if _, exists := totals[sortedKeys[k]]; exists {
			last = totals[sortedKeys[k]]
			change = last - results[sortedKeys[k]]
		}

		output = output + fmt.Sprintf("%03d_%03d_%+04d_%s\n\r", results[sortedKeys[k]], last, change, sortedKeys[k])
		newTotals[sortedKeys[k]] = results[sortedKeys[k]]
	}
	totals = newTotals
	history = append(history, newTotals)

	sendResults(fmt.Sprintf("Error Counts, Every %smin", freq), output)
}

func sendResults(t, s string) {
	mstClient := goteamsnotify.NewClient()

	card := goteamsnotify.NewMessageCard()
	card.Title = t
	card.Text = s

	if err := mstClient.Send(webhookUrl, card); err != nil {
		log.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}
}
