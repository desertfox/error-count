package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/desertfox/gograylog"
)

func doQuery() string {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		return string(data)
	}

	c := gograylog.New(os.Getenv("EC_HOST"), os.Getenv("EC_USER"), os.Getenv("EC_PASS"))

	f, _ := strconv.Atoi(freq)

	data, err := c.Execute(os.Getenv("EC_QUERY"), os.Getenv("EC_STREAMID"), []string{"message"}, 10000, f)
	if err != nil {
		fmt.Println("Unable to make graylog request", err)
		return ""
	}
	return string(data)
}
