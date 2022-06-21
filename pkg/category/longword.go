package category

import (
	"context"
	"regexp"
	"sort"
	"strings"
)

var (
	magic         *regexp.Regexp
	r             *regexp.Regexp
	AlphaNumColon string = "alphanumcolon"
)

func init() {
	r, _ = regexp.Compile("Logger")
	magic, _ = regexp.Compile("[^a-zA-Z:/ -]+")
}

func CreateKeyFn() func(ctx context.Context, b []byte) (string, error) {
	return func(ctx context.Context, b []byte) (string, error) {
		return CreateKey(b), nil
	}
}

func CreateKey(b []byte) string {
	words := strings.Split(magic.ReplaceAllString(string(b), ""), " ")

	sort.SliceStable(words, func(i int, j int) bool {
		return len(words[i]) < len(words[j])
	})

	var largest string
	for i := len(words); i > 0; i-- {
		if r.MatchString(words[i-1]) {
			continue
		} else if len(words[i-1]) > 120 {
			continue
		}
		largest = words[i-1]
		break
	}

	return largest
}
