package power

import (
	"../hnap"
)

type HNAP struct {
	Control
}

func (control HNAP) On()  {
	log.Info("HNAP On event fired")
	hnap.On()
}

func (control HNAP) Off()  {
	log.Info("HNAP Off event fired")
	hnap.Off()
}

func (control HNAP) State() bool {
	log.Info("HNAP state event fired")
	return hnap.State()
}

func (control HNAP) Consumption() float64 {
	log.Info("HNAP consumption event fired")
	return hnap.Consumption()
}
