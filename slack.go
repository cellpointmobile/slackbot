package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"strings"
	"./event"
	"os"
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
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)

				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					respond(rtm, ev, prefix)
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				//Take no action
			}
		}
	}
}

func respond(rtm *slack.RTM, msg *slack.MessageEvent, prefix string) {
	var response string
	text := msg.Text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	turnStuffOn := map[string]bool{
		"coffee aye?": true,
		"COFFEE NOW!":       true,
		"I'm dying here":         true,
		"on": true,
	}
	turnStuffOff := map[string]bool{
		"please stop": true,
		"ditch the black liquid":     true,
		"off":   true,
	}

	if turnStuffOn[text] {
		response = "okay okay, relax dude.."
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
		event.Power_on()
	} else if turnStuffOff[text] {
		response = "Terminating coffee supplies!"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
		event.Power_off()
	}
}