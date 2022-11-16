package main

import (
	"bytes"
	"text/template"

	"desertfox.dev/error-count/v1/pkg/count"
)

var (
	//header string = "<pre>DAY  HOUR PREV NOW  +/- FIRST_SEEN FILE <br>"
	headers []string = []string{"DAY", "HOUR", "PREV", "NOW", "+/-", "FIRST_SEEN", "FILE"}
	//line       string = "%04d_%04d_%04d_%04d_%+05d %s %s %d <br>"
	timeFormat string = "2006-01-02 15:04:05"
)

/*
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
*/

type line struct {
	day, hour, prev, now, diff int
	seen, file                 string
	line                       int
}

func totals(day, hour, prev, last count.Ledger) string {
	t, _ := template.New("teams").Parse(`
	<table>
		<tr>
			{{range .Headers }}
				<th>{{.}}</th>
			{{end}}
		</tr>\n
		{{ range .Data }}
			<tr>
				<td>{{.day}}</td>
				<td>{{.hour}}</td>
				<td>{{.prev}}</td>
				<td>{{.now}}</td>
				<td>{{.diff}}</td>
				<td>{{.seen}}</td>
				<td>{{.file}}:{{.line}}</td>
			</tr>\n
		{{end}}
	</table>`)

	var lines []line
	for _, file := range day.GetTopFileInstances(30) {
		d := day.GetCount(file)
		h := hour.GetCount(file)
		p := prev.GetCount(file)
		c := last.GetCount(file)

		lines = append(lines, line{d.Count, h.Count, p.Count, c.Count, c.Count - p.Count, count.TimeLedger[file].Format(timeFormat), d.Record.File, d.Record.Line})
	}

	var b bytes.Buffer
	err := t.Execute(&b, struct {
		Headers []string
		Data    []line
	}{
		Headers: headers,
		Data:    lines,
	})

	if err != nil {
		panic(err)
	}

	return b.String()
}
