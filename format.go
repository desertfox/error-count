package main

import (
	"fmt"

	"desertfox.dev/error-count/v1/pkg/count"
)

var (
	header     string = "<pre>DAY  HOUR PREV NOW  +/- FIRST_SEEN FILE <br>"
	line       string = "%04d_%04d_%04d_%04d_%+05d %s %s %d <br>"
	timeFormat string = "2006-01-02 15:04:05"
)

func totals(day, hour, prev, last count.Ledger) string {
	var output string = header
	for _, file := range day.GetTopFileInstances(30) {
		d := day.GetCount(file)
		h := hour.GetCount(file)
		p := prev.GetCount(file)
		c := last.GetCount(file)

		output = output + fmt.Sprintf(line, d.Count, h.Count, p.Count, c.Count, c.Count-p.Count, count.TimeLedger[file].Format(timeFormat), d.Record.File, d.Record.Line)
	}
	output = output + "</pre>"

	return output
}
