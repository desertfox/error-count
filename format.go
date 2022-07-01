package main

import (
	"fmt"

	"desertfox.dev/error-count/v1/pkg/count"
)

var (
	header string = "COUNT_PREV_+/-_FILE <br> "
)

func totals(ls count.Ledgers) string {
	l := ls.GetLast()
	pL := ls.GetPrev()

	var output string = header
	for _, file := range l.GetTopFileInstances(10) {
		c := l.GetCount(file)
		pc := pL.GetCount(file)

		output = output + fmt.Sprintf("* %03d_%03d_%+04d_%s:%d <br> ", c.Count, pc.Count, c.Count-pc.Count, c.Record.File, c.Record.Line)
	}

	return output
}
