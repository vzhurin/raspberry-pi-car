package chassis

import (
	"errors"
	"math"
	"time"
)

type pin interface {
	Out(bool) error
}

type pwm interface {
	Start(dutyCycle float64, frequency float64, errChan chan<- error) error
	Stop()
}

type Chassis struct {
	pwmControllerRight pwm
	pwmControllerLeft  pwm

	controlPinRight1 pin
	controlPinRight2 pin

	controlPinLeft1 pin
	controlPinLeft2 pin

	pwmFrequency float64
	moveDuration time.Duration
}

func NewChassis(
	pwmControllerRight pwm,
	pwmControllerLeft pwm,
	controlPinRight1 pin,
	controlPinRight2 pin,
	controlPinLeft1 pin,
	controlPinLeft2 pin,
	pwmFrequency float64,
	moveDuration time.Duration,
) *Chassis {
	return &Chassis{
		pwmControllerRight: pwmControllerRight,
		pwmControllerLeft:  pwmControllerLeft,
		controlPinRight1:   controlPinRight1,
		controlPinRight2:   controlPinRight2,
		controlPinLeft1:    controlPinLeft1,
		controlPinLeft2:    controlPinLeft2,
		pwmFrequency:       pwmFrequency,
		moveDuration:       moveDuration,
	}
}

func (c *Chassis) Move(rightDutyCycle float64, leftDutyCycle float64) error {
	if rightDutyCycle < -1 || rightDutyCycle > 1 || leftDutyCycle < -1 || leftDutyCycle > 1 {
		return errors.New("rightDutyCycle and leftDutyCycle values must be in the range from -1 to 1 inclusive")
	}

	return c.move(rightDutyCycle, leftDutyCycle, c.moveDuration)
}

func (c *Chassis) move(rightDutyCycle float64, leftDutyCycle float64, duration time.Duration) error {
	if rightDutyCycle > 0 {
		err := c.controlPinRight1.Out(false)
		if err != nil {
			return err
		}

		err = c.controlPinRight2.Out(true)
		if err != nil {
			return err
		}
	} else {
		err := c.controlPinRight1.Out(true)
		if err != nil {
			return err
		}

		err = c.controlPinRight2.Out(false)
		if err != nil {
			return err
		}
	}

	if leftDutyCycle > 0 {
		err := c.controlPinLeft1.Out(true)
		if err != nil {
			return err
		}

		err = c.controlPinLeft2.Out(false)
		if err != nil {
			return err
		}
	} else {
		err := c.controlPinLeft1.Out(false)
		if err != nil {
			return err
		}

		err = c.controlPinLeft2.Out(true)
		if err != nil {
			return err
		}
	}

	rightErrChan := make(chan error)
	err := c.pwmControllerRight.Start(math.Abs(rightDutyCycle), c.pwmFrequency, rightErrChan)
	if err != nil {
		return err
	}

	leftErrChan := make(chan error)
	err = c.pwmControllerLeft.Start(math.Abs(leftDutyCycle), c.pwmFrequency, leftErrChan)
	if err != nil {
		return err
	}

	timer := time.NewTimer(duration)

	select {
	case err := <-rightErrChan:
		return err
	case err := <-leftErrChan:
		return err
	case <-timer.C:
	}

	c.pwmControllerRight.Stop()
	c.pwmControllerLeft.Stop()

	return nil
}
