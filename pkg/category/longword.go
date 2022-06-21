package category

import (
	"context"
	"regexp"
	"sort"
	"strings"
)

var (
	reg           *regexp.Regexp
	AlphaNumColon string = "alphanumcolon"
)

func init() {
	reg, _ = regexp.Compile("[^a-zA-Z:/ _-]+")
}

func CreateKeyFn() func(ctx context.Context, b []byte) (string, error) {
	return func(ctx context.Context, b []byte) (string, error) {
		return CreateKey(b), nil
	}
}

func CreateKey(b []byte) string {
	words := strings.Split(reg.ReplaceAllString(string(b), " "), " ")

	sort.SliceStable(words, func(i int, j int) bool {
		return len(words[i]) < len(words[j])
	})

	return words[len(words)-1]
}
