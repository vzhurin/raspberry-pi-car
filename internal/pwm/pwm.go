package pwm

import (
	"errors"
	"fmt"
	"raspberry-pi-car/internal/pin"
	"time"
)

type PWM struct {
	pin  pin.Pin
	done chan struct{}
}

func NewPWM(pin pin.Pin) *PWM {
	return &PWM{
		pin:  pin,
		done: make(chan struct{}),
	}
}

func (p *PWM) Start(dutyCycle uint, frequency uint) error {
	err := p.validate(dutyCycle, frequency)
	if err != nil {
		return err
	}

	period := time.Second / time.Duration(frequency)

	go p.work(period, dutyCycle)

	return nil
}

func (p *PWM) Stop() {
	p.done <- struct{}{}
}

func (p *PWM) work(period time.Duration, dutyCycle uint) {
	highDuration := (period * time.Duration(dutyCycle)) / time.Duration(100)
	lowDuration := period - highDuration

	defer func() {
		_ = p.pin.Out(pin.Low)
	}()

	edgeCase := false
	level := pin.Low
	if dutyCycle == 0 {
		edgeCase = true
		level = pin.Low
	} else if dutyCycle == 100 {
		edgeCase = true
		level = pin.High
	}

	if edgeCase {
		_ = p.pin.Out(level)
		fmt.Println("DONE")
		<-p.done
		fmt.Println("DONE DONE")

		return
	}

	for {
		_ = p.pin.Out(pin.High)
		time.Sleep(highDuration)

		_ = p.pin.Out(pin.Low)
		time.Sleep(lowDuration)

		select {
		case <-p.done:
			return
		default:
		}
	}
}

func (p *PWM) validate(dutyCycle uint, frequency uint) error {
	if dutyCycle > 100 {
		return errors.New("duty cycle must not be more than 100")
	}

	if frequency == 0 {
		return errors.New("frequency must not be zero")
	}

	if time.Second/time.Duration(frequency) == 0 {
		return errors.New("frequency is too high")
	}

	return nil
}
