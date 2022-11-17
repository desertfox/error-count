package main

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	fileRe     = regexp.MustCompile(`/[/\w+-]+`)
	libFileRe  = regexp.MustCompile(`/lib/[/\w+-]+`)
	filelineRe = regexp.MustCompile(`line \d+`)
)

func FileLineKeyFn() func(ctx context.Context, s string) (string, int, error) {
	return func(ctx context.Context, s string) (string, int, error) {
		return FileLine(s)
	}
}

func FileLine(s string) (string, int, error) {
	if filelineRe.MatchString(s) {
		return getFile(s), getLine(s), nil
	}
	return "", 0, errors.New("no match found")
}

func getFile(s string) string {
	f := libFileRe.FindString(s)
	parts := strings.Split(f, "/")

	if len(parts) == 0 {
		return fileRe.FindString(s)
	}

	if len(parts) < 2 {
		fmt.Println(s, f)
		return strings.Join(parts, "::")
	}

	var chomp int = 2
	if len(parts) > 1 && parts[2] == "perl5" {
		chomp = 3
	}

	return strings.Join(parts[chomp:], "::")
}

func getLine(s string) int {
	lineat := filelineRe.FindString(s)
	line, _ := strconv.Atoi(strings.Split(lineat, " ")[1])
	return line
}
