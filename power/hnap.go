package power

import (
	"fmt"
	"../hnap"
)

type HNAP struct {
	Control
}

func (control HNAP) On()  {
	fmt.Println("HNAP On event fired")
	hnap.On()
}

func (control HNAP) Off()  {
	fmt.Println("HNAP Off event fired")
	hnap.Off()
}

func (control HNAP) State() bool {
	fmt.Println("HNAP consumption event fired")
	return hnap.State()
}

func (control HNAP) Consumption() float64 {
	fmt.Println("HNAP consumption event fired")
	return hnap.Consumption()
}
