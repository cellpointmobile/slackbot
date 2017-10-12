package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"strings"
	"./event"
	"./power"
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
		"coffee now!":       true,
		"i'm dying here":         true,
		"on": true,
	}
	turnStuffOff := map[string]bool{
		"please stop": true,
		"ditch the black liquid":     true,
		"off":   true,
	}

	randomResponses := map[string]string{
		"die!": "Never!",
		"milk?": "Ask <@U1JCMNPHC> to go get it..",
	}

	impl := new (power.HNAP)

	if turnStuffOn[text] {
		response = "okay okay, relax dude.."
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
		impl.On()
	} else if turnStuffOff[text] {
		response = "Terminating coffee supplies!"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
		impl.Off()
	} else if randomResponses[text] != "" {
		rtm.SendMessage(rtm.NewOutgoingMessage(randomResponses[text], msg.Channel))
	} else {
		author, quote := event.Get_random_quote()
		rtm.SendMessage(rtm.NewOutgoingMessage(quote + "\n\n - _" + author + "_", msg.Channel))
	}
}