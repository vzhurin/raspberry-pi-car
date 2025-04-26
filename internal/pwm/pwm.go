package pwm

import (
	"errors"
	"sync"
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
	pin     pin
	started bool
	done    chan struct{}
	wg      *sync.WaitGroup
}

func NewPWM(pin pin) *PWM {
	return &PWM{
		pin:     pin,
		started: false,
		done:    make(chan struct{}),
		wg:      new(sync.WaitGroup),
	}
}

func (p *PWM) Start(dutyCycle uint, frequency uint) error {
	err := p.validate(dutyCycle, frequency)
	if err != nil {
		return err
	}

	p.Stop()

	if dutyCycle == 100 {
		err = p.pin.Out(High)
		if err != nil {
			return err
		}

		return nil
	}

	if dutyCycle == 0 {
		err = p.pin.Out(Low)
		if err != nil {
			return err
		}

		return nil
	}

	period := time.Second / time.Duration(frequency)
	p.work(period, dutyCycle)

	return nil
}

func (p *PWM) Stop() {
	if p.started {
		p.done <- struct{}{}
	}

	p.wg.Wait()
}

func (p *PWM) work(period time.Duration, dutyCycle uint) {
	ticker := time.NewTicker(period)

	p.wg.Add(1)
	p.started = true
	go func() {
		defer func() {
			ticker.Stop()
			_ = p.pin.Halt()
			p.started = false
			p.wg.Done()
		}()
		for {
			select {
			case <-ticker.C:
				_ = p.pin.Out(High)
				time.Sleep((period * time.Duration(dutyCycle)) / time.Duration(100))
				_ = p.pin.Out(Low)
			case <-p.done:
				return
			}
		}
	}()
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
