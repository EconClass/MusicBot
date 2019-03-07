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

/*
   TODO: Change @BOT_NAME to the same thing you entered when creating your Slack application.
   NOTE: command_arg_1 and command_arg_2 represent optional parameteras that you define
   in the Slack API UI
*/
const helpMessage = `*COMMANDS:*
"@LyricBot *Query*"
>*Query* MUST BE the name of an artist, album, playlist or track.`

//CreateSlackClient sets up the slack RTM (real-timemessaging) client library,
//initiating the socket connection and returning the client.
func CreateSlackClient(apiKey string) *slack.RTM {
	api := slack.New(apiKey)
	rtm := api.NewRTM()
	go rtm.ManageConnection() // goroutine!
	return rtm
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
			case message == "cats":
				sendCats(slackClient, message, ev.Channel)
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

//Cat does stuff
type Cat struct {
	Text string `json:"text"`
}

func sendCats(slackClient *slack.RTM, message, slackChannel string) {
	command := strings.ToLower(message)
	url := "https://cat-fact.herokuapp.com/facts/random"
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if err != nil {
		log.Fatal(getErr)
	}

	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	cat := &Cat{}
	json.Unmarshal([]byte(body), &cat)
	fmt.Printf("%+v\n", cat)
	// bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// bodyString := string(bodyBytes)

	fmt.Println("[RECEIVED] sendCats:", command)
	slackClient.SendMessage(slackClient.NewOutgoingMessage(cat.Text, slackChannel))
}

type Lyrics struct {
	LyricsText string `json:"lyrics_body"`
}

type LyricsObj struct {
	LyrObj Lyrics `json:"lyrics"`
}

type MusixBody struct {
	MessageBody LyricsObj `json:"body"`
}

type Musix struct {
	Message MusixBody `json:"message"`
}

// sendResponse is NOT unimplemented --- write code in the function body to complete!
func sendResponse(slackClient *slack.RTM, message, slackChannel string) {
	musicKey := os.Getenv("KEY")
	command := strings.ToLower(message)
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
	}

	fmt.Printf("%+v\n", musix)
	// bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// bodyString := string(bodyBytes)

	fmt.Println("[RECEIVED] sendResponse:", command)
	slackClient.SendMessage(slackClient.NewOutgoingMessage(musix.Message.MessageBody.LyrObj.LyricsText, slackChannel))
}
