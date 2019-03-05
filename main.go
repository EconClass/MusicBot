package main

import (
	"net/http"
	"os"

	musicbot "github.com/EconClass/MusicBot/slackstuff"

	_ "github.com/joho/godotenv/autoload"
)

// main is our entrypoint, where the application initializes the Slackbot.
func main() {
	port := ":" + os.Getenv("PORT")
	go http.ListenAndServe(port, nil)
	slackIt()
}

// slackIt is a function that initializes the Slackbot.
func slackIt() {
	botToken := os.Getenv("BOT_OAUTH_ACCESS_TOKEN")
	slackClient := musicbot.CreateSlackClient(botToken)
	musicbot.RespondToEvents(slackClient)
}
