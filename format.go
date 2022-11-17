package main

import (
	"bytes"
	"text/template"
)

var (
	headers    []string = []string{"DAY", "HOUR", "PREV", "NOW", "+/-", "FIRST_SEEN", "FILE"}
	timeFormat string   = "2006-01-02 15:04:05"
)

type line struct {
	Day, Hour, Prev, Now, Diff int
	Seen, File                 string
	Line                       int
}

func totals(day, hour, prev, last Ledger) string {
	t, _ := template.New("teams").Parse(`<table><tr>{{range .Headers }}<th>{{.}}</th>{{end}}</tr>{{ range .Data }}<tr><td>{{.Day}}</td><td>{{.Hour}}</td><td>{{.Prev}}</td><td>{{.Now}}</td><td>{{.Diff}}</td><td>{{.Seen}}</td><td>{{.File}}:{{.Line}}</td></tr>{{end}}</table>`)

	var lines []line
	for _, file := range day.GetTopFileInstances(30) {
		d := day.GetCount(file)
		h := hour.GetCount(file)
		p := prev.GetCount(file)
		c := last.GetCount(file)

		lines = append(lines, line{d.Count, h.Count, p.Count, c.Count, c.Count - p.Count, TimeLedger[file].Format(timeFormat), d.Record.File, d.Record.Line})
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
