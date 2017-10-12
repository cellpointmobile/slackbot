package power

type Control interface {
	On()
	Off()
	State() bool
	Consumption() float64
}
