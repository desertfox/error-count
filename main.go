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
	hLedgers          = make(count.Ledgers, 0)
	dLedgers          = make(count.Ledgers, 0)
)

func main() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(freq + "m").Do(doInterval)

	s.Every("60m").Do(doHour)

	s.Every("24h").Do(doDay)

	s.StartBlocking()
}

func doInterval() {
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
	hLedgers.Add(ledger)

	teams.SendResults(
		webhookUrl,
		fmt.Sprintf("%sm Error Counts", freq),
		totals(hLedgers),
	)
}

func doHour() {
	ledger := hLedgers.TotalLedger()
	dLedgers.Add(ledger)

	teams.SendResults(webhookUrl, "1h Error Count.", totals(dLedgers))

	hLedgers = make(count.Ledgers, 0)
}

func doDay() {
	ledger := dLedgers.TotalLedger()

	teams.SendResults(webhookUrl, "24h Error Count.", totals(count.Ledgers{ledger}))

	dLedgers = make(count.Ledgers, 0)
}
