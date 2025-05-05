package pwm

import (
	"errors"
	"raspberry-pi-car/internal/pin"
	"sync"
	"time"
)

type PWM struct {
	pin  pin.Pin
	done chan struct{}
	wg   sync.WaitGroup
}

func NewPWM(pin pin.Pin) *PWM {
	return &PWM{
		pin:  pin,
		done: make(chan struct{}),
		wg:   sync.WaitGroup{},
	}
}

func (p *PWM) Start(dutyCycle float64, frequency float64) error {
	err := p.validate(dutyCycle, frequency)
	if err != nil {
		return err
	}

	period := time.Duration(float64(time.Second) / frequency)
	highDuration := time.Duration(float64(period) * dutyCycle)
	lowDuration := period - highDuration

	p.wg.Add(1)
	go p.work(highDuration, lowDuration)

	return nil
}

func (p *PWM) Stop() {
	p.done <- struct{}{}
	p.wg.Wait()
}

func (p *PWM) work(highDuration, lowDuration time.Duration) {
	defer func() {
		_ = p.pin.Out(pin.Low)
		p.wg.Done()
	}()

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

func (p *PWM) validate(dutyCycle float64, frequency float64) error {
	if dutyCycle < 0 || dutyCycle > 1 {
		return errors.New("duty cycle must be in the range from 0 to 1 inclusive")
	}

	if frequency <= 0 {
		return errors.New("frequency must be greater than zero")
	}

	return nil
}
