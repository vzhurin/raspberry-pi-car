package pwm

import (
	"errors"
	"sync"
	"time"
)

type pin interface {
	Out(bool) error
}

type PWM struct {
	pin  pin
	done chan struct{}
	wg   sync.WaitGroup
}

func NewPWM(pin pin) *PWM {
	return &PWM{
		pin:  pin,
		done: make(chan struct{}),
		wg:   sync.WaitGroup{},
	}
}

func (p *PWM) Start(dutyCycle float64, frequency float64, errChan chan<- error) error {
	err := p.validate(dutyCycle, frequency)
	if err != nil {
		return err
	}

	period := time.Duration(float64(time.Second) / frequency)
	highDuration := time.Duration(float64(period) * dutyCycle)
	lowDuration := period - highDuration

	p.wg.Add(1)
	go p.work(highDuration, lowDuration, errChan)

	return nil
}

func (p *PWM) Stop() {
	p.done <- struct{}{}
	p.wg.Wait()
}

func (p *PWM) work(highDuration, lowDuration time.Duration, errChan chan<- error) {
	defer func() {
		_ = p.pin.Out(false)
		p.wg.Done()
	}()

	for {
		err := p.pin.Out(true)
		if err != nil {
			errChan <- err

			return
		}
		time.Sleep(highDuration)

		err = p.pin.Out(false)
		if err != nil {
			errChan <- err

			return
		}
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
