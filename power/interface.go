package power

import "github.com/op/go-logging"

type Control interface {
	On()
	Off()
	State() bool
	Consumption() float64
}

var log = logging.MustGetLogger("")
