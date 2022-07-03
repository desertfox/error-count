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
	freq            string = os.Getenv("EC_FREQ")
	webhookUrl      string = os.Getenv("EC_TEAMSWEBHOOK")
	intervalLedgers        = make(count.Ledgers, 0)
	hourLedgers            = make(count.Ledgers, 0)
	start                  = time.Now()
)

func main() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(freq + "m").Do(doInterval)

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
	intervalLedgers.Add(ledger)

	if len(intervalLedgers) == 6 {
		hourLedgers.Add(intervalLedgers.TotalLedger())
	}

	var t count.Ledgers = hourLedgers
	if len(hourLedgers) == 0 {
		t = count.Ledgers{intervalLedgers.TotalLedger()}
	}

	teams.SendResults(
		webhookUrl,
		fmt.Sprintf("%sm Error Counts. Uptime:%s", freq, time.Since(start)),
		totals(
			t.TotalLedger(),
			t.GetLast(),
			intervalLedgers.GetPrev(),
			intervalLedgers.GetLast(),
		),
	)

	if len(intervalLedgers) == 6 {
		intervalLedgers = make(count.Ledgers, 0)
	}

	if len(hourLedgers) == 24 {
		hourLedgers = make(count.Ledgers, 0)
	}

}
