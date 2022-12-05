package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/desertfox/gograylog"
	"github.com/go-co-op/gocron"
)

var (
	start, reset             time.Time = time.Now(), start.Add(time.Hour * 24)
	poolCount, resetInterval int       = 5, 60 / 15
)

func main() {
	var (
		timeLedger                   TimeLedger        = make(TimeLedger, 0)
		hourLedgers, intervalLedgers Ledgers           = make(Ledgers, 0), make(Ledgers, 0)
		s                            *gocron.Scheduler = gocron.NewScheduler(time.UTC)
	)

	s.Every(freq+"m").Do(interval, timeLedger, &hourLedgers, &intervalLedgers, query)

	s.StartBlocking()
}

func query() []string {
	gg := gograylog.New(graylogHost, graylogUser, os.Getenv("EC_PASS"))

	f, _ := strconv.Atoi(freq)

	data, err := gg.Execute(graylogQuery, graylogStreamID, []string{"message"}, 10000, f)
	if err != nil {
		fmt.Println("Unable to make graylog request", err)
		return []string{}
	}
	return strings.Split(string(data), "\n")
}

func notify(t, s string) {
	mstClient := goteamsnotify.NewClient()
	mstClient.SkipWebhookURLValidationOnSend(true)

	card := goteamsnotify.NewMessageCard()
	card.Title, card.Text = t, s

	if err := mstClient.Send(webhookUrl, card); err != nil {
		log.Printf(
			"failed to send message: %v",
			err,
		)
	}
}

func interval(timeLedger TimeLedger, hourLedger, intervalLedger *Ledgers, query func() []string) {
	var (
		jobs   []Job
		ledger Ledger = make(Ledger, 0)
	)

	for _, line := range query() {
		jobs = append(jobs, Job{
			Data: line,
			Fnc:  fileLineFnc(),
		})
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	wp := Pool{
		instances: poolCount,
		jobs:      make(chan Job, poolCount),
		results:   make(chan Record, poolCount),
	}
	go wp.Queue(jobs)
	go wp.Run(ctx)

	for record := range wp.Results() {
		if record.Err != nil {
			fmt.Println(record.Err)
			continue
		}
		ledger.Update(record)
		timeLedger.Track(record)
	}
	intervalLedger.Add(ledger)

	if len(*intervalLedger) >= resetInterval {
		hourLedger.Add(intervalLedger.Total())

		defer func() {
			*intervalLedger = make(Ledgers, 0)
		}()
	}

	hL := *hourLedger
	switch len(*hourLedger) {
	case 0:
		hL = Ledgers{intervalLedger.Total()}
	case 24:
		defer func() {
			*hourLedger = make(Ledgers, 0)
		}()
	}

	title := fmt.Sprintf(
		"error-count-%s, report every %sm. uptime %.fh totals reset in %.fh",
		teamsTitle,
		freq,
		time.Since(start).Hours(),
		time.Until(reset).Hours(),
	)
	totals := total(hL.Total(), hL.Last(), intervalLedger.Prev(), intervalLedger.Last(), timeLedger)
	notify(title, totals)

	if time.Now().After(reset) {
		reset = time.Now().Add(time.Hour * 24)
	}
}
