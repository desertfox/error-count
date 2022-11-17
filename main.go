package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
)

var (
	intervalLedgers           = make(Ledgers, 0)
	hourLedgers               = make(Ledgers, 0)
	start                     = time.Now()
	reset           time.Time = start.Add(time.Hour * 24)
)

func main() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(freq + "m").Do(doInterval)

	s.StartBlocking()
}

func doInterval() {
	lines := strings.Split(doQuery(), "\n")
	jobs := make([]Job, len(lines))

	for i := range lines {
		jobs[i] = Job{
			Data:   lines[i],
			ExecFn: FileLineKeyFn(),
		}
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	wp := NewWP(5)
	go wp.Queue(jobs)
	go wp.Run(ctx)

	ledger := NewLedger()
	for r := range wp.Results() {
		if r.Err == nil {
			ledger.Add(r)
		}
	}
	intervalLedgers.Add(ledger)

	if len(intervalLedgers) == 6 {
		hourLedgers.Add(intervalLedgers.TotalLedger())
	}

	var t Ledgers = hourLedgers
	if len(hourLedgers) == 0 {
		t = Ledgers{intervalLedgers.TotalLedger()}
	}

	SendResults(
		webhookUrl,
		fmt.Sprintf("error-count-%s, report every %sm. uptime %.fh totals reset in %.fh", teamsTitle, freq, time.Since(start).Hours(), time.Until(reset).Hours()),
		totals(
			t.TotalLedger(),
			t.GetLast(),
			intervalLedgers.GetPrev(),
			intervalLedgers.GetLast(),
		),
	)

	if len(intervalLedgers) == 6 {
		intervalLedgers = make(Ledgers, 0)
	}

	if len(hourLedgers) == 24 {
		hourLedgers = make(Ledgers, 0)
	}

	if time.Now().After(reset) {
		reset = time.Now().Add(time.Hour * 24)
	}

}
