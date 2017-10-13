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
	"sync"
	"math/rand"
)

var
(
	log = logging.MustGetLogger("")
	rtm *slack.RTM
	channel string = "C7GNULHMK"
	state bool
	stateMutex sync.Mutex
	inProgressMsg = [4]string{"Still brewing..", "Just a few more jiffies..", "Control your caffeine urge! I'm not done yet..", "ZzzZzzZ! Slowly pressing the brown juice.."}
)

func main() {

	level, err := logging.LogLevel(os.Getenv("LOG_LEVEL") )
	if err == nil {
		botlogging.SetupLogging(level)
	} else {
		botlogging.SetupLogging(logging.WARNING)
	}

	token := os.Getenv("SLACK_TOKEN")
	if len(token) < 1  {
		log.Fatal("Please set environment variable SLACK_TOKEN")
	}
	setChannel := os.Getenv("SLACK_CHANNEL")
	if len(setChannel) > 0  {
		channel = setChannel
	} else {
		log.Debug("SLACK_CHANNEL not set in environment, using default: " + channel)
	}

	api := slack.New(token)
	rtm = api.NewRTM()

	go coffee_thread()
	slack_thread()
}

func slack_thread() {
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
	state = impl.State()

	for {
		time.Sleep(3 * time.Second)
		stateMutex.Lock()
		newState := impl.State()
		if state != newState {

			msg := ""
			if newState {
				msg = "Let there be brew! Some marvellous angel decided to cook a batch of heavenly caffeine for us"
			} else {
				msg = "Dang! One naughty soul foiled our brewing plot"
			}

			rtm.SendMessage(rtm.NewOutgoingMessage(msg, channel))
			state = newState
		}
		stateMutex.Unlock()
	}
}

func brew_coffee(rtm *slack.RTM, msg *slack.MessageEvent) {
	var response string
	impl := new (power.HNAP)

	stateMutex.Lock()
		impl.On()
	state = true
	stateMutex.Unlock()

	Loop:
		for {
			time.Sleep(30 * time.Second)

			if impl.Consumption() >= 100 {
				response = inProgressMsg[rand.Intn(4)]
			} else {
				response = "Brew completed! :coffee: :tada:"
				break Loop
			}
			log.Debug("Sending message: " + response + " to channel: " + msg.Channel)
			rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
		}

	rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
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
		"coffee please": true,
		"caffeine running low": true,
		"on": true,
	}
	turnStuffOff := map[string]bool{
		"please stop": true,
		"caffeine overflow": true,
		"i'm chuck norris": true,
		"off": true,
	}

	randomResponses := map[string]string{
		"die!": "Never!",
		"milk?": "Ask <@U1JCMNPHC> to go get it..",
	}

	impl := new (power.HNAP)

	if turnStuffOn[text] {
		response := ""
		if state == false {
			response = "okay okay, relax dude.."
			go brew_coffee(rtm, msg)
		} else if (impl.Consumption() >= 100) {
			response = inProgressMsg[rand.Intn(4)]
		} else {
			response = "Pffff.. Yesterdays news ya landlobber! How 'bout sloppin' what's already on da pot?!"
		}
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))

	} else if turnStuffOff[text] {
		response = "Terminating coffee supplies!"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))

		stateMutex.Lock()
			impl.Off()
		state = false
		stateMutex.Unlock()

	} else if randomResponses[text] != "juice usage?" {
		energy := impl.Consumption()
		juice := "Could not read juice level :pensive:"
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