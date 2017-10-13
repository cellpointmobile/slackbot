package main

import (
	"os"
	"./power"
	"fmt"
	"strconv"
)

func main() {
	execute(os.Args[1])
}

func execute(command string) {
	// TODO: Make driver switchable from the outside
	impl := new (power.HNAP)

	if command == "on" {
		impl.On()
	} else if command == "off" {
		impl.Off()
	} else if command == "state" {
		state := impl.State()
		if state {
			fmt.Println("State: on")
		} else {
			fmt.Println("State: off")
		}
	} else if command == "consumption" {
		energy := impl.Consumption()
		if energy >= 0 {
			fmt.Println(strconv.FormatFloat(energy, 'f', 2, 64) + " watts")
		}
	} else if command == "demo" {
		execute("on")
		execute("state")
		execute("consumption")
		execute("off")
		execute("state")
		execute("consumption")
	} else {
		fmt.Println("Unknown command: " + command)
	}
}
