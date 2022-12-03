package main

import (
	"bytes"
	"text/template"
)

var (
	headers []string = []string{"DAY", "HOUR", "PREV", "NOW", "+/-", "FIRST_SEEN", "FILE"}
	format  string   = "2006-01-02 15:04:05"
)

type tableRow struct {
	Day, Hour, Prev, Now, Diff int
	Seen, File                 string
	Line                       int
}

func total(day, hour, prev, last Ledger, tL TimeLedger) string {
	t, _ := template.New("teams").Parse(`<table><tr>{{range .Headers }}<th>{{.}}</th>{{end}}</tr>{{ range .Data }}<tr><td>{{.Day}}</td><td>{{.Hour}}</td><td>{{.Prev}}</td><td>{{.Now}}</td><td>{{.Diff}}</td><td>{{.Seen}}</td><td>{{.File}}:{{.Line}}</td></tr>{{end}}</table>`)

	var (
		rows []tableRow
		b    bytes.Buffer
	)
	for _, file := range day.Top(30) {
		d := day.Get(file)
		h := hour.Get(file)
		p := prev.Get(file)
		c := last.Get(file)

		rows = append(rows, tableRow{d.Count, h.Count, p.Count, c.Count, c.Count - p.Count, tL[file].Format(format), d.Record.File, d.Record.Line})
	}

	err := t.Execute(&b, struct {
		Headers []string
		Data    []tableRow
	}{
		Headers: headers,
		Data:    rows,
	})

	if err != nil {
		panic(err)
	}

	return b.String()
}
