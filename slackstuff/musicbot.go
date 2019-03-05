package musicbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/nlopes/slack"
)

/*
   TODO: Change @BOT_NAME to the same thing you entered when creating your Slack application.
   NOTE: command_arg_1 and command_arg_2 represent optional parameteras that you define
   in the Slack API UI
*/
const helpMessage = `*COMMANDS:*
"@spots *Query*"
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

			// TODO: Make your bot do more than respond to a help command. See notes below.
			// Make changes below this line and add additional funcs to support your bot's functionality.
			// sendHelp is provided as a simple example. Your team may want to call a free external API
			// in a function called sendResponse that you'd create below the definition of sendHelp,
			// and call in this context to ensure execution when the bot receives an event.

			// START SLACKBOT CUSTOM CODE
			// ===============================================================
			switch {
			case message == "help":
				sendHelp(slackClient, message, ev.Channel)
			default:
				sendCats(slackClient, message, ev.Channel)
			}
			// ===============================================================
			// END SLACKBOT CUSTOM CODE
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

// sendResponse is NOT unimplemented --- write code in the function body to complete!
func sendResponse(slackClient *slack.RTM, message, slackChannel string) {
	command := strings.ToLower(message)

	// START SLACKBOT CUSTOM CODE
	// ===============================================================
	// TODO:
	//      1. Implement sendResponse for one or more of your custom Slackbot commands.
	//         You could call an external API here, or create your own string response. Anything goes!
	//      2. STRETCH: Write a goroutine that calls an external API based on the data received in this function.
	// ===============================================================
	// END SLACKBOT CUSTOM CODE
	fmt.Println("[RECEIVED] sendResponse:", command)
}
