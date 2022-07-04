package category

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	fileRe     = regexp.MustCompile(`/lib/[/\w+-]+`)
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
	return "unknown", 0, errors.New("no match found")
}

func getFile(s string) string {
	parts := strings.Split(fileRe.FindString(s), "/")

	if len(parts) < 2 {
		return strings.Join(parts, "::")
	}

	return strings.Join(parts[2:], "::")
}

func getLine(s string) int {
	lineat := filelineRe.FindString(s)
	line, _ := strconv.Atoi(strings.Split(lineat, " ")[1])
	return line
}
