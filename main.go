package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/desertfox/gograylog"
	"github.com/go-co-op/gocron"
)

var start, reset time.Time = time.Now(), start.Add(time.Hour * 24)

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
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		return strings.Split(string(data), "\n")
	}

	c := gograylog.New(graylogHost, graylogUser, os.Getenv("EC_PASS"))

	f, _ := strconv.Atoi(freq)

	data, err := c.Execute(graylogQuery, graylogStreamID, []string{"message"}, 10000, f)
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

func interval(totalLedger TimeLedger, hourLedger, intervalLedger *Ledgers, query func() []string) {
	rawLines := query()
	jobs := make([]Job, len(rawLines))

	for i := range rawLines {
		jobs[i] = Job{
			Data: rawLines[i],
			Fnc:  fileLineFnc(),
		}
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	wp := Pool{
		instances: 5,
		jobs:      make(chan Job, 5),
		results:   make(chan Record, 5),
	}
	go wp.Queue(jobs)
	go wp.Run(ctx)

	ledger := make(Ledger, 0)
	for record := range wp.Results() {
		if record.Err != nil {
			fmt.Println(record.Err)
		}
		ledger.Incriment(record)
		totalLedger.Add(record)
	}
	intervalLedger.Add(ledger)

	if len(*intervalLedger) == 6 {
		hourLedger.Add(intervalLedger.Total())
	}

	var t Ledgers = *hourLedger
	if len(*hourLedger) == 0 {
		t = Ledgers{intervalLedger.Total()}
	}

	title := fmt.Sprintf(
		"error-count-%s, report every %sm. uptime %.fh totals reset in %.fh",
		teamsTitle,
		freq,
		time.Since(start).Hours(),
		time.Until(reset).Hours(),
	)
	totals := total(t.Total(), t.Last(), intervalLedger.Prev(), intervalLedger.Last(), totalLedger)
	notify(title, totals)

	if len(*intervalLedger) == 6 {
		*intervalLedger = make(Ledgers, 0)
	}

	if len(*hourLedger) == 24 {
		*hourLedger = make(Ledgers, 0)
	}

	if time.Now().After(reset) {
		reset = time.Now().Add(time.Hour * 24)
	}
}
