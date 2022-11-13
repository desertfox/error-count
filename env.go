package main

import "os"

var (
	freq        string = os.Getenv("EC_FREQ")
	webhookUrl  string = os.Getenv("EC_TEAMSWEBHOOK")
	graylogHost string = os.Getenv("EC_HOST")
	graylogUser string = os.Getenv("EC_USER")
	//password left off
	graylogQuery    string = os.Getenv("EC_QUERY")
	graylogStreamID string = os.Getenv("EC_STREAMID")
	teamsTitle      string = os.Getenv("EC_TITLE")
)
