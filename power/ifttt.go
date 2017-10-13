package power

import (
	"net/http"
	"os"
)

const IFTTT_BASE_URL = "https://maker.ifttt.com/trigger"

type IFTTT struct {
	Control
}

func (control IFTTT) On()  {
	log.Info("IFTTT On event fired")
	http.Get(IFTTT_BASE_URL + "/power_on/with/key/" + os.Getenv("IFTTT_TOKEN") )
}

func (control IFTTT) Off()  {
	log.Info("IFTTT Off event fired")
	http.Get(IFTTT_BASE_URL + "/power_off/with/key/" + os.Getenv("IFTTT_TOKEN") )
}

func (control IFTTT) State() bool {
	log.Info("IFTTT State not available")
	return false
}

func (control IFTTT) Consumption() float64 {
	log.Info("IFTTT Consumption not available")
	return -1
}
