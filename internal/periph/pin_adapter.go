package periph

import (
	"periph.io/x/conn/v3/gpio"
	"raspberry-pi-car/internal/pwm"
)

type Pin struct {
	pin gpio.PinOut
}

func NewPin(pin gpio.PinOut) *Pin {
	return &Pin{
		pin: pin,
	}
}

func (p *Pin) Out(level pwm.Level) error {
	if level == pwm.High {
		return p.pin.Out(gpio.High)
	}

	return p.pin.Out(gpio.Low)
}

func (p *Pin) Halt() error {
	return p.pin.Halt()
}
