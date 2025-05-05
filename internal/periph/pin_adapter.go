package periph

import (
	"periph.io/x/conn/v3/gpio"
)

type Pin struct {
	pin gpio.PinOut
}

func NewPin(pin gpio.PinOut) *Pin {
	return &Pin{
		pin: pin,
	}
}

func (p *Pin) Out(level bool) error {
	if level == true {
		return p.pin.Out(gpio.High)
	}

	return p.pin.Out(gpio.Low)
}
