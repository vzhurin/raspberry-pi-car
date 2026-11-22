package chassis

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"
)

type pin interface {
	Out(bool) error
}

type pwm interface {
	Start(ctx context.Context, dutyCycle float64, freqHz float64) (<-chan error, error)
	Stop()
}

type Chassis struct {
	mu                 sync.Mutex
	pwmControllerRight pwm
	pwmControllerLeft  pwm

	controlPinRight1 pin
	controlPinRight2 pin

	controlPinLeft1 pin
	controlPinLeft2 pin

	pwmFreqHz    float64
	moveDuration time.Duration
}

func New(
	pwmControllerRight pwm,
	pwmControllerLeft pwm,
	controlPinRight1 pin,
	controlPinRight2 pin,
	controlPinLeft1 pin,
	controlPinLeft2 pin,
	pwmFreqHz float64,
	moveDuration time.Duration,
) *Chassis {
	return &Chassis{
		pwmControllerRight: pwmControllerRight,
		pwmControllerLeft:  pwmControllerLeft,
		controlPinRight1:   controlPinRight1,
		controlPinRight2:   controlPinRight2,
		controlPinLeft1:    controlPinLeft1,
		controlPinLeft2:    controlPinLeft2,
		pwmFreqHz:          pwmFreqHz,
		moveDuration:       moveDuration,
	}
}

func (c *Chassis) Move(ctx context.Context, rightGas, leftGas float64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := validateGas(rightGas); err != nil {
		return err
	}

	if err := validateGas(leftGas); err != nil {
		return err
	}

	if rightGas > 0 {
		if err := c.setRightForward(); err != nil {
			return err
		}
	} else {
		if err := c.setRightBackward(); err != nil {
			return err
		}
	}

	if leftGas > 0 {
		if err := c.setLeftForward(); err != nil {
			return err
		}
	} else {
		if err := c.setLeftBackward(); err != nil {
			return err
		}
	}

	moveCtx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		c.pwmControllerRight.Stop()
		c.pwmControllerLeft.Stop()
	}()

	rightErrChan, err := c.pwmControllerRight.Start(
		moveCtx,
		math.Abs(rightGas),
		c.pwmFreqHz,
	)
	if err != nil {
		return err
	}

	leftErrChan, err := c.pwmControllerLeft.Start(
		moveCtx,
		math.Abs(leftGas),
		c.pwmFreqHz,
	)
	if err != nil {
		return err
	}

	timer := time.NewTimer(c.moveDuration)
	defer timer.Stop()

	select {
	case err := <-rightErrChan:
		return err
	case err := <-leftErrChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (c *Chassis) setRightForward() error {
	err := c.controlPinRight1.Out(true)
	if err != nil {
		return err
	}

	err = c.controlPinRight2.Out(false)
	if err != nil {
		return err
	}

	return nil
}

func (c *Chassis) setRightBackward() error {
	err := c.controlPinRight1.Out(false)
	if err != nil {
		return err
	}

	err = c.controlPinRight2.Out(true)
	if err != nil {
		return err
	}

	return nil
}

func (c *Chassis) setLeftForward() error {
	err := c.controlPinLeft1.Out(false)
	if err != nil {
		return err
	}

	err = c.controlPinLeft2.Out(true)
	if err != nil {
		return err
	}

	return nil
}

func (c *Chassis) setLeftBackward() error {
	err := c.controlPinLeft1.Out(true)
	if err != nil {
		return err
	}

	err = c.controlPinLeft2.Out(false)
	if err != nil {
		return err
	}

	return nil
}

func validateGas(gas float64) error {
	if gas < -1 || gas > 1 {
		return errors.New("invalid gas")
	}

	return nil
}
