package chassis

import (
	"errors"
	"math"
	"raspberry-pi-car/internal/pin"
	"raspberry-pi-car/internal/pwm"
	"time"
)

const pwmFreq = 100

type Chassis struct {
	pwmRight *pwm.PWM
	pwmLeft  *pwm.PWM

	controlPinRight1 pin.Pin
	controlPinRight2 pin.Pin

	controlPinLeft1 pin.Pin
	controlPinLeft2 pin.Pin
}

func NewChassis(
	pwmRight *pwm.PWM,
	pwmLeft *pwm.PWM,
	controlPinRight1 pin.Pin,
	controlPinRight2 pin.Pin,
	controlPinLeft1 pin.Pin,
	controlPinLeft2 pin.Pin,
) *Chassis {
	return &Chassis{
		pwmRight:         pwmRight,
		pwmLeft:          pwmLeft,
		controlPinRight1: controlPinRight1,
		controlPinRight2: controlPinRight2,
		controlPinLeft1:  controlPinLeft1,
		controlPinLeft2:  controlPinLeft2,
	}
}

func (c *Chassis) Move(rightDutyCycle float64, leftDutyCycle float64, duration time.Duration) error {
	if rightDutyCycle < -1 || rightDutyCycle > 1 || leftDutyCycle < -1 || leftDutyCycle > 1 {
		return errors.New("rightDutyCycle and leftDutyCycle values must be in the range from -1 to 1 inclusive")
	}

	return c.move(rightDutyCycle, leftDutyCycle, duration)
}

func (c *Chassis) move(rightDutyCycle float64, leftDutyCycle float64, duration time.Duration) error {
	if rightDutyCycle > 0 {
		err := c.controlPinRight1.Out(pin.Low)
		if err != nil {
			return err
		}

		err = c.controlPinRight2.Out(pin.High)
		if err != nil {
			return err
		}
	} else {
		err := c.controlPinRight1.Out(pin.High)
		if err != nil {
			return err
		}

		err = c.controlPinRight2.Out(pin.Low)
		if err != nil {
			return err
		}
	}

	if leftDutyCycle > 0 {
		err := c.controlPinLeft1.Out(pin.High)
		if err != nil {
			return err
		}

		err = c.controlPinLeft2.Out(pin.Low)
		if err != nil {
			return err
		}
	} else {
		err := c.controlPinLeft1.Out(pin.Low)
		if err != nil {
			return err
		}

		err = c.controlPinLeft2.Out(pin.High)
		if err != nil {
			return err
		}
	}

	err := c.pwmRight.Start(math.Abs(rightDutyCycle), pwmFreq)
	if err != nil {
		return err
	}

	err = c.pwmLeft.Start(math.Abs(leftDutyCycle), pwmFreq)
	if err != nil {
		return err
	}

	time.Sleep(duration)

	c.pwmRight.Stop()
	c.pwmLeft.Stop()

	return nil
}
