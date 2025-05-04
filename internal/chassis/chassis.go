package chassis

import (
	"errors"
	"raspberry-pi-car/internal/helper"
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

func (c *Chassis) Move(left int, right int, duration time.Duration) error {
	if left < -100 || right < -100 || left > 100 || right > 100 {
		return errors.New("left and right values must be in the range from -100 to 100 inclusive")
	}

	return c.move(left, right, duration)
}

func (c *Chassis) move(left int, right int, duration time.Duration) error {
	if right > 0 {
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

	if left > 0 {
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

	rightWorker := func(errChan chan<- error) {
		errChan <- c.pwmRight.Start(uint(helper.Abs(right)), pwmFreq)
	}

	leftWorker := func(errChan chan<- error) {
		errChan <- c.pwmLeft.Start(uint(helper.Abs(left)), pwmFreq)
	}

	errChan := make(chan error)

	go rightWorker(errChan)
	go leftWorker(errChan)

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	time.Sleep(duration)

	c.pwmRight.Stop()
	c.pwmLeft.Stop()

	return nil
}
