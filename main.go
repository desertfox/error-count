package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"desertfox.dev/error-count/v1/pkg/category"
	"desertfox.dev/error-count/v1/pkg/count"
	"desertfox.dev/error-count/v1/pkg/teams"
	"desertfox.dev/error-count/v1/pkg/worker"
	"github.com/go-co-op/gocron"
)

var (
	freq       string = os.Getenv("EC_FREQ")
	webhookUrl string = os.Getenv("EC_TEAMSWEBHOOK")
	ledgers           = make(count.Ledgers, 0)
)

func main() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(freq + "m").Do(do)
	s.Every("60m").Do(report)
	s.StartBlocking()
}

func do() {
	lines := strings.Split(doQuery(), "\n")
	jobs := make([]worker.Job, len(lines))

	for i := range lines {
		jobs[i] = worker.Job{
			Data:   lines[i],
			ExecFn: category.FileLineKeyFn(),
		}
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	wp := worker.NewWP(5)
	go wp.Queue(jobs)
	go wp.Run(ctx)

	ledger := count.NewLedger()
	for r := range wp.Results() {
		if r.Err == nil {
			ledger.Add(r)
		}
	}
	ledgers.Add(ledger)

	output := "COUNT_PREV_+/-_FILE\n\r"
	prevLedger := ledgers.GetLast()
	for _, file := range ledger.GetTopFileInstances(10) {
		c := ledger.GetCount(file)
		pc := prevLedger.GetCount(file)

		output = output + fmt.Sprintf("%03d_%03d_%+04d_%s:%d\n\r", c.Count, pc.Count, c.Count-pc.Count, c.Record.File, c.Record.Line)
	}

	teams.SendResults(webhookUrl, fmt.Sprintf("Error Counts, Every %smin", freq), output)
}

func report() {
	ledger := ledgers.TotalLedger()
	output := "COUNT_FILE\n\r"
	for _, file := range ledger.GetTopFileInstances(10) {
		c := ledger.GetCount(file)

		output = output + fmt.Sprintf("%03d_%s\n\r", c.Count, c.Record.File)
	}
	teams.SendResults(webhookUrl, "Error Count Hour Totals", output)

	if len(ledgers) >= 6 {
		ledgers = make(count.Ledgers, 0)
	}
}
