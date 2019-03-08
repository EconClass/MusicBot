package musicbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

const helpMessage = `*COMMANDS:*
"@LyricBot *Query*"
>*Query* MUST BE the name of a track.`

//CreateSlackClient sets up the slack RTM (real-timemessaging) client library,
//initiating the socket connection and returning the client.
func CreateSlackClient(apiKey string) *slack.RTM {
	api := slack.New(apiKey)
	rtm := api.NewRTM()
	go rtm.ManageConnection() // goroutine!
	return rtm
}

//Lyrics is used to model the json respose from Musixmatch API
type Lyrics struct {
	LyricsText string `json:"lyrics_body"`
}

//LyricsObj is used to model the json respose from Musixmatch API
type LyricsObj struct {
	LyrObj Lyrics `json:"lyrics"`
}

//MusixBody is used to model the json respose from Musixmatch API
type MusixBody struct {
	MessageBody LyricsObj `json:"body"`
}

//Musix is used to model the json respose from Musixmatch API
type Musix struct {
	Message MusixBody `json:"message"`
}

//RespondToEvents waits for messages on the Slack client's incomingEvents channel,
//and sends a response when it detects the bot has been tagged in a message with @<botTag>.
func RespondToEvents(slackClient *slack.RTM) {
	for msg := range slackClient.IncomingEvents {
		fmt.Println("Event Received: ", msg.Type)
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			botTagString := fmt.Sprintf("<@%s> ", slackClient.GetInfo().User.ID)
			if !strings.Contains(ev.Msg.Text, botTagString) {
				continue
			}
			message := strings.Replace(ev.Msg.Text, botTagString, "", -1)

			switch {
			case message == "help":
				sendHelp(slackClient, message, ev.Channel)
			default:
				sendResponse(slackClient, message, ev.Channel)
			}
		default:
		}
	}
}

// sendHelp is a working help message, for reference.
func sendHelp(slackClient *slack.RTM, message, slackChannel string) {
	if strings.ToLower(message) != "help" {
		return
	}

	slackClient.SendMessage(slackClient.NewOutgoingMessage(helpMessage, slackChannel))
}

func sendResponse(slackClient *slack.RTM, message, slackChannel string) {
	musicKey := os.Getenv("KEY")
	command := strings.ToLower(message)
	command = strings.Replace(command, " ", "%20", -1)
	fmt.Println(command)
	baseURL := "https://api.musixmatch.com/ws/1.1/matcher.lyrics.get?format=json&callback=callback&q_track="
	url := baseURL + command + "&apikey=" + musicKey
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if err != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	musix := &Musix{}
	jsonErr := json.Unmarshal([]byte(body), &musix)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		slackClient.SendMessage(slackClient.NewOutgoingMessage("Try a different song.", slackChannel))
	}

	// fmt.Printf("%+v\n", musix)

	fmt.Println("[RECEIVED] sendResponse:", command)
	slackClient.SendMessage(slackClient.NewOutgoingMessage(musix.Message.MessageBody.LyrObj.LyricsText, slackChannel))
}
