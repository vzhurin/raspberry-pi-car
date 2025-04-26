package pwm

import (
	"errors"
	"fmt"
	"time"
)

type pin interface {
	Out(level Level) error
	Halt() error
}

type Level bool

const (
	High Level = true
	Low  Level = false
)

type PWM struct {
	pin  pin
	done chan struct{}
}

func NewPWM(pin pin) *PWM {
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

	return p.work(period, dutyCycle)
}

func (p *PWM) Stop() {
	p.done <- struct{}{}
}

func (p *PWM) work(period time.Duration, dutyCycle uint) (err error) {
	highDuration := (period * time.Duration(dutyCycle)) / time.Duration(100)
	lowDuration := period - highDuration

	defer func() {
		err = p.pin.Out(Low)
		err = p.pin.Halt()
		fmt.Println("PWM stopped")
	}()

	level := Low
	edgeCase := false
	if dutyCycle == 0 {
		edgeCase = true
		level = Low

	}

	if dutyCycle == 100 {
		edgeCase = true
		level = High
	}

	if edgeCase {
		err := p.pin.Out(level)
		if err != nil {
			return err
		}

		<-p.done

		return nil
	}

	for {
		err := p.pin.Out(High)
		if err != nil {
			return err
		}
		time.Sleep(highDuration)

		err = p.pin.Out(Low)
		if err != nil {
			return err
		}
		time.Sleep(lowDuration)

		select {
		case <-p.done:
			return nil
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
