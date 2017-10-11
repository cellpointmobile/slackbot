package event

import (
	"fmt"
	"net/http"
	"os"
)

const IFTTT_BASE_URL = "https://maker.ifttt.com/trigger"

func Power_on() {
	fmt.Println("power on event fired")
	http.Get(IFTTT_BASE_URL + "/power_on/with/key/" + os.Getenv("IFTTT_TOKEN") )
}

func Power_off() {
	fmt.Println("power off event fired")
	http.Get(IFTTT_BASE_URL + "/power_off/with/key/" + os.Getenv("IFTTT_TOKEN") )
}
