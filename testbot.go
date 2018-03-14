package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

func main() {

	token := os.Getenv("SLACK_TOKEN")
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:

	for {
		select {

		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event received: ")

			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s>", info.User.ID)

				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					rtm.SendMessage(rtm.NewOutgoingMessage("What's up buddy?", ev.Channel))
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				// No action; do Go switches have to be exhaustive?

			}
		}
	}
}

func respond(rtm *slack.RTM, msg *slack.MessageEvent, prefix string) {
	var response string // I guess the zero value is used, so this syntax is used instead
	text := msg.Text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	acceptedGreetings := map[string]bool{
		"what's up?": true,
		"hey!":       true,
		"yo":         true,
	}
	acceptedHowAreYou := map[string]bool{
		"how's it going?": true,
		"how are ya?":     true,
		"feeling okay?":   true,
	}

	if acceptedGreetings[text] { //Uses the fact that the zero value of Bool is false?
		response = "What's up buddy?"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	} else if acceptedHowAreYou[text] {
		response = "Good, how are you?"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	}
}
