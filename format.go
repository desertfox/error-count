package main

import (
	"fmt"

	"desertfox.dev/error-count/v1/pkg/count"
)

var (
	header string = "DAY HOUR PREV NOW +/- FILE <br>"
	line   string = "%4d_%4d_%4d_%4d_%+5d--->%s:%d <br>"
)

func totals(day, hour, prev, last count.Ledger) string {
	var output string = header
	for _, file := range day.GetTopFileInstances(15) {
		d := day.GetCount(file)
		h := hour.GetCount(file)
		p := prev.GetCount(file)
		c := last.GetCount(file)

		output = output + fmt.Sprintf(line, d.Count, h.Count, p.Count, c.Count, c.Count-p.Count, d.Record.File, d.Record.Line)
	}

	return output
}
