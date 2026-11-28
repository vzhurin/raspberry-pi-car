package app

import (
	"log/slog"
	"os"
	"raspberry-pi-car/config"
	"raspberry-pi-car/internal/chassis"
	"raspberry-pi-car/internal/hardware"
	"raspberry-pi-car/internal/joystick"
	"raspberry-pi-car/internal/joystick/0xcafed00d"
	goSDL3 "raspberry-pi-car/internal/joystick/go-sdl3"
	"raspberry-pi-car/internal/pin"
	"raspberry-pi-car/internal/pwm"

	"periph.io/x/conn/v3/gpio/gpioreg"
)

type Container struct {
	cfg *config.Config

	pwmControllerRight *pwm.PWM
	pwmControllerLeft  *pwm.PWM

	chassis *chassis.Chassis

	watcher *hardware.Watcher

	joystickDispatcher *joystick.Dispatcher
	joystickListener   *joystick.Listener

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

		c.pwmControllerRight = pwm.New(pin.NewPin(pwmPinRight))
	}

	return c.pwmControllerRight
}

func (c *Container) PWMControllerLeft() *pwm.PWM {
	if c.pwmControllerLeft == nil {
		pwmPinLeft := gpioreg.ByName(c.cfg.PWMPinLeft)
		if pwmPinLeft == nil {
			panic("failed to find pwm pin left")
		}

		c.pwmControllerLeft = pwm.New(pin.NewPin(pwmPinLeft))
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
			pin.NewPin(controlPinRight1),
			pin.NewPin(controlPinRight2),
			pin.NewPin(controlPinLeft1),
			pin.NewPin(controlPinLeft2),
			c.cfg.PWMFrequency,
			c.cfg.MoveDuration,
		)
	}

	return c.chassis
}

func (c *Container) Watcher() *hardware.Watcher {
	if c.watcher == nil {
		c.watcher = hardware.NewWatcher(
			c.Logger().With(slog.String("where", "watcher")),
			c.cfg.Joystick,
		)
	}

	return c.watcher
}

func (c *Container) JoystickDispatcher() *joystick.Dispatcher {
	if c.joystickDispatcher == nil {
		c.joystickDispatcher = joystick.NewDispatcher(
			c.Logger().With(slog.String("where", "dispatcher")),
			c.JoystickGoSDL3,
			c.Watcher().Bus(),
		)
	}

	return c.joystickDispatcher
}

func (c *Container) JoystickListener() *joystick.Listener {
	if c.joystickListener == nil {
		c.joystickListener = joystick.NewListener(
			c.Chassis(),
			c.Logger().With(slog.String("where", "listener")),
			c.JoystickDispatcher().StateBus(),
		)
	}

	return c.joystickListener
}

func (c *Container) Joystick0xcafed00d(device string) (joystick.Joystick, error) {
	id, err := hardware.DeviceID(device)
	if err != nil {
		return nil, err
	}

	return _0xcafed00d.New(
		c.Logger().With(slog.String("where", "joystick 0xcafed00d")),
		id,
		c.cfg.MoveDuration,
	)
}

func (c *Container) JoystickGoSDL3(device string) (joystick.Joystick, error) {
	id, err := hardware.DeviceID(device)
	if err != nil {
		return nil, err
	}

	return goSDL3.New(
		c.Logger().With(slog.String("where", "joystick go SDL3")),
		id,
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
