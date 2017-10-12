package power

import (
	"net/http"
	"os"
	"fmt"
)

const IFTTT_BASE_URL = "https://maker.ifttt.com/trigger"

type IFTTT struct {
	Control
}

func (control IFTTT) On()  {
	fmt.Println("IFTTT On event fired")
	http.Get(IFTTT_BASE_URL + "/power_on/with/key/" + os.Getenv("IFTTT_TOKEN") )
}

func (control IFTTT) Off()  {
	fmt.Println("IFTTT Off event fired")
	http.Get(IFTTT_BASE_URL + "/power_off/with/key/" + os.Getenv("IFTTT_TOKEN") )
}

func (control IFTTT) State() bool {
	fmt.Println("IFTTT State not available")
	return false
}

func (control IFTTT) Consumption() float64 {
	fmt.Println("IFTTT Consumption not available")
	return -1
}
