package app

import (
	"log/slog"
	"os"
	"raspberry-pi-car/config"
	"raspberry-pi-car/internal/chassis"
	"raspberry-pi-car/internal/joystick"
	"raspberry-pi-car/internal/periph"
	"raspberry-pi-car/internal/pwm"

	"periph.io/x/conn/v3/gpio/gpioreg"
)

type Container struct {
	cfg *config.Config

	pwmControllerRight *pwm.PWM
	pwmControllerLeft  *pwm.PWM

	chassis *chassis.Chassis

	watcher    *joystick.Watcher
	dispatcher *joystick.Dispatcher
	joystick   *joystick.Joystick

	logger *slog.Logger
}

func NewContainer(cfg *config.Config) *Container {
	return &Container{cfg: cfg}
}

func (c *Container) PWMControllerRight() *pwm.PWM {
	if c.pwmControllerRight == nil {
		pwmPinRight := gpioreg.ByName(c.cfg.PWMPinRight)
		if pwmPinRight == nil {
			panic("failed to find pwm pin right")
		}

		c.pwmControllerRight = pwm.New(periph.NewPin(pwmPinRight))
	}

	return c.pwmControllerRight
}

func (c *Container) PWMControllerLeft() *pwm.PWM {
	if c.pwmControllerLeft == nil {
		pwmPinLeft := gpioreg.ByName(c.cfg.PWMPinLeft)
		if pwmPinLeft == nil {
			panic("failed to find pwm pin left")
		}

		c.pwmControllerLeft = pwm.New(periph.NewPin(pwmPinLeft))
	}

	return c.pwmControllerLeft
}

func (c *Container) Chassis() *chassis.Chassis {
	if c.chassis == nil {
		controlPinRight1 := gpioreg.ByName(c.cfg.ControlPinRight1)
		if controlPinRight1 == nil {
			panic("failed to find right control pin 1")
		}

		controlPinRight2 := gpioreg.ByName(c.cfg.ControlPinRight2)
		if controlPinRight2 == nil {
			panic("failed to find right control pin 2")
		}

		controlPinLeft1 := gpioreg.ByName(c.cfg.ControlPinLeft1)
		if controlPinLeft1 == nil {
			panic("failed to find left control pin 1")
		}

		controlPinLeft2 := gpioreg.ByName(c.cfg.ControlPinLeft2)
		if controlPinLeft2 == nil {
			panic("failed to find left control pin 2")
		}

		c.chassis = chassis.New(
			c.PWMControllerRight(),
			c.PWMControllerLeft(),
			periph.NewPin(controlPinRight1),
			periph.NewPin(controlPinRight2),
			periph.NewPin(controlPinLeft1),
			periph.NewPin(controlPinLeft2),
			c.cfg.PWMFrequency,
			c.cfg.MoveDuration,
		)
	}

	return c.chassis
}

func (c *Container) Watcher() *joystick.Watcher {
	if c.watcher == nil {
		c.watcher = joystick.NewWatcher(
			c.Logger().With(slog.String("where", "watcher")),
			make(chan joystick.Event, 1),
			c.cfg.JoystickID,
		)
	}

	return c.watcher
}

func (c *Container) Dispatcher() *joystick.Dispatcher {
	if c.dispatcher == nil {
		c.dispatcher = joystick.NewDispatcher(
			c.Logger().With(slog.String("where", "dispatcher")),
			c.Joystick,
			c.Watcher().Bus(),
		)
	}

	return c.dispatcher
}

func (c *Container) Joystick(joystickID int) (*joystick.Joystick, error) {
	return joystick.New(
		c.Chassis(),
		c.Logger().With(slog.String("where", "joystick")),
		joystickID,
		c.cfg.MoveDuration,
	)
}

func (c *Container) Logger() *slog.Logger {
	if c.logger == nil {
		c.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	return c.logger
}
