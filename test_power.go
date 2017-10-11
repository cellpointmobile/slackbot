package main

import (
	"os"
	"./event"
	"fmt"
)

func main() {

	command := os.Args[1]

	if command == "on" {
		event.Power_on()
	} else if command == "off" {
		event.Power_off()
	} else {
		fmt.Println("Unknown command: " + command)
	}
}
