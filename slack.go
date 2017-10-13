package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"strings"
	"./event"
	"./power"
	"os"
	"strconv"
	"time"
	"github.com/op/go-logging"
	"./botlogging"
)

var log = logging.MustGetLogger("")

func main() {

	level, err := logging.LogLevel(os.Getenv("LOG_LEVEL") )
	if err == nil {
		botlogging.SetupLogging(level)
	} else {
		botlogging.SetupLogging(logging.WARNING)
	}

	token := os.Getenv("SLACK_TOKEN")
	if os.Getenv("SLACK_CHANNEL") == "" {
		os.Setenv("SLACK_CHANNEL", "#hackerdaysdk")
	}
	go coffee_thread()
	slack_thread(token)
}

func slack_thread(token string) {
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			log.Debug("Event Received")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				log.Info("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				log.Info("Message: %v\n", ev)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)

				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					respond(rtm, ev, prefix)
				}

			case *slack.RTMError:
				log.Error(fmt.Sprintf("Error: %s\n", ev.Error()))

			case *slack.InvalidAuthEvent:
				log.Error("Invalid credentials")
				break Loop

			default:
				//Take no action
			}
		}
	}
}

func coffee_thread() {
	// get initial status
	impl := new (power.HNAP)
	state := impl.State()

	for {
		time.Sleep(3 * time.Second)
		if state != impl.State() {
			log.Debug("State changed to: " + strconv.FormatBool(impl.State()))
		}
	}
}

func brew_coffee(rtm *slack.RTM, msg *slack.MessageEvent) {
	var response string
	impl := new (power.HNAP)
	impl.On()

	Loop:
		for {
			if impl.Consumption() >= 100 {
				response = "Still brewing.."
			} else {
				response = "Brew completed!"
				break Loop
			}
			log.Debug("Sending message: " + response + " to channel: " + msg.Channel)
			rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
			time.Sleep(30 * time.Second)
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
		"coffee now!": true,
		"i'm dying here": true,
		"on": true,
	}
	turnStuffOff := map[string]bool{
		"please stop": true,
		"ditch the black liquid": true,
		"off": true,
	}

	randomResponses := map[string]string{
		"die!": "Never!",
		"milk?": "Ask <@U1JCMNPHC> to go get it..",
	}

	impl := new (power.HNAP)

	if turnStuffOn[text] {
		brew_coffee(rtm, msg)
		response := "okay okay, relax dude.."
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	} else if turnStuffOff[text] {
		response = "Terminating coffee supplies!"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
		impl.Off()
	} else if randomResponses[text] != "juice usage?" {
		energy := impl.Consumption()
		juice := "Could not read juice level :("
		if energy >= 0 {
			juice = strconv.FormatFloat(energy, 'f', 2, 64) + " watts"
		}
		rtm.SendMessage(rtm.NewOutgoingMessage(juice, msg.Channel))
	} else if randomResponses[text] != "" {
		rtm.SendMessage(rtm.NewOutgoingMessage(randomResponses[text], msg.Channel))
	} else {
		author, quote := event.Get_random_quote()
		rtm.SendMessage(rtm.NewOutgoingMessage(quote + "\n\n - _" + author + "_", msg.Channel))
	}
}