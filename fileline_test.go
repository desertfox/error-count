package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func Test_fileline(t *testing.T) {
	data, err := ioutil.ReadFile("./sample_data/gl.csv")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")

	for i := range lines {
		r := line(lines[i]).parse()
		t.Log(fmt.Sprintf("%s:%s:%d", lines[i], r.File, r.Location))
	}

}
