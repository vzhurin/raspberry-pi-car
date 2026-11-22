package pwm

import (
	"context"
	"errors"
	"sync"
	"time"
)

type pin interface {
	Out(bool) error
}

type PWM struct {
	pin    pin
	mu     sync.Mutex
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func New(pin pin) *PWM {
	return &PWM{
		pin: pin,
	}
}

func (p *PWM) Start(ctx context.Context, dutyCycle float64, freqHz float64) (<-chan error, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cancel != nil {
		return nil, errors.New("PWM already started")
	}

	if err := validateDutyCycle(dutyCycle); err != nil {
		return nil, err
	}
	if err := validateFrequency(freqHz); err != nil {
		return nil, err
	}

	period := time.Duration(float64(time.Second) / freqHz)
	highDuration := time.Duration(float64(period) * dutyCycle)

	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel

	errChan := make(chan error, 1)

	p.wg.Add(1)
	go p.work(ctx, period, highDuration, errChan)

	return errChan, nil
}

func (p *PWM) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cancel == nil {
		return
	}

	p.cancel()
	p.wg.Wait()
	p.cancel = nil
}

func (p *PWM) work(ctx context.Context, period, highDuration time.Duration, errChan chan<- error) {
	defer func() {
		_ = p.pin.Out(false)
		p.wg.Done()
	}()

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if highDuration > 0 {
				if err := p.pin.Out(true); err != nil {
					errChan <- err
					return
				}
			}
			if highDuration < period {
				timer := time.NewTimer(highDuration)

				select {
				case <-ctx.Done():
					timer.Stop()
					return
				case <-timer.C:
				}

				if err := p.pin.Out(false); err != nil {
					errChan <- err
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func validateDutyCycle(dutyCycle float64) error {
	if dutyCycle < 0 || dutyCycle > 1 {
		return errors.New("invalid duty cycle")
	}

	return nil
}

func validateFrequency(frequency float64) error {
	if frequency <= 0 {
		return errors.New("invalid frequency")
	}

	return nil
}
