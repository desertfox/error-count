package category

import (
	"context"
	"regexp"
)

var (
	filelineRe = regexp.MustCompile(`(/[/\w+-. ]+) line \d+`)
)

func FileLineKeyFn() func(ctx context.Context, b []byte) (string, error) {
	return func(ctx context.Context, b []byte) (string, error) {
		return FileLine(b)
	}
}

func FileLine(b []byte) (string, error) {
	if filelineRe.Match(b) {
		file := filelineRe.FindString(string(b))

		//fmt.Printf("raw\n%s\nfile\n%s\n\n", b, file)

		return file, nil
	}
	return "unknown", nil
}
