// This project utilizes the starter kit provided by @droxey:
// https://github.com/droxey/goslackit
package main

import (
	"net/http"
	"os"

	musicbot "musicbot/slackstuff"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	port := ":" + os.Getenv("PORT")
	go http.ListenAndServe(port, nil)
	slackIt()
}

func slackIt() {
	botToken := os.Getenv("BOT_OAUTH_ACCESS_TOKEN")
	slackClient := musicbot.CreateSlackClient(botToken)
	musicbot.RespondToEvents(slackClient)
}
