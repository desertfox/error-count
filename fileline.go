package main

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var regex map[string]*regexp.Regexp = map[string]*regexp.Regexp{
	"file":     regexp.MustCompile(`/[/\w+-]+`),
	"lib":      regexp.MustCompile(`/lib/[/\w+-]+`),
	"location": regexp.MustCompile(`line \d+`),
}

type line string

func fileLineFnc() func(ctx context.Context, s string) Record {
	return func(ctx context.Context, s string) Record {
		return line(s).parse()
	}
}

func (l line) parse() Record {
	if l.hasLocation() {
		return Record{
			File:     l.file(),
			Location: l.location(),
		}
	}
	return Record{
		Err: errors.New("no match found"),
	}
}

func (l line) hasLocation() bool {
	return regex["location"].MatchString(string(l))
}

func (l line) file() string {
	f := regex["lib"].FindString(string(l))

	parts := strings.Split(f, "/")

	if len(parts) == 0 {
		return regex["file"].FindString(string(l))
	}

	if len(parts) < 2 {
		return strings.Join(parts, "::")
	}

	var chomp int = 2
	if len(parts) > 1 && parts[2] == "perl5" {
		chomp = 3
	}

	return strings.Join(parts[chomp:], "::")
}

func (l line) location() int {
	lineat := strings.Split(regex["location"].FindString(string(l)), " ")[1]
	line, _ := strconv.Atoi(lineat)
	return line
}
