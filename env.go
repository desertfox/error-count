package main

import "os"

var (
	freq                                                    string = os.Getenv("EC_FREQ")
	webhookUrl, teamsTitle                                  string = os.Getenv("EC_TEAMSWEBHOOK"), os.Getenv("EC_TITLE")
	graylogHost, graylogUser, graylogQuery, graylogStreamID string = os.Getenv("EC_HOST"), os.Getenv("EC_USER"), os.Getenv("EC_QUERY"), os.Getenv("EC_STREAMID")
)
