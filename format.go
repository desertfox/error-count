package main

import (
	"fmt"

	"desertfox.dev/error-count/v1/pkg/count"
)

var (
	header string = "DAY_HOUR_PREV_NOW_+/-_FILE <br>"
	line   string = "%03d_%03d_%03d_%03d_%+04d_%s:%d <br>"
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
