package app

import (
	"log/slog"
	"os"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"raspberry-pi-car/config"
	apiHTTP "raspberry-pi-car/internal/api/http"
	"raspberry-pi-car/internal/chassis"
	"raspberry-pi-car/internal/periph"
	"raspberry-pi-car/internal/pwm"
	"sync"
	"time"
)

type Container struct {
	cfg          *config.Config
	httpServerMu sync.Mutex
	chassisMu    sync.Mutex

	httpServer  *apiHTTP.Server
	httpHandler *apiHTTP.Handler

	pwmControllerRight *pwm.PWM
	pwmControllerLeft  *pwm.PWM

	chassis *chassis.Chassis

	logger *slog.Logger
}

func NewContainer(cfg *config.Config) *Container {
	return &Container{cfg: cfg}
}

func (c *Container) HTTPServer() *apiHTTP.Server {
	c.httpServerMu.Lock()
	defer c.httpServerMu.Unlock()

	if c.httpServer == nil {
		httpServer, err := apiHTTP.NewServer(&apiHTTP.Config{
			Host:    c.cfg.RESTHost,
			Port:    c.cfg.RESTPort,
			GinMode: c.cfg.GinMode,
		}, c.HTTPHandler(), c.Logger())

		if err != nil {
			panic(err)
		}

		c.httpServer = httpServer
	}

	return c.httpServer
}

func (c *Container) HTTPHandler() *apiHTTP.Handler {
	if c.httpHandler == nil {
		c.httpHandler = apiHTTP.NewHandler(c.Chassis(), c.Logger())
	}

	return c.httpHandler
}

func (c *Container) PWMControllerRight() *pwm.PWM {
	if c.pwmControllerRight == nil {
		pwmPinRight := gpioreg.ByName(c.cfg.PWMPinRight)
		if pwmPinRight == nil {
			panic("failed to find pwm pin right")
		}

		c.pwmControllerRight = pwm.NewPWM(periph.NewPin(pwmPinRight))
	}

	return c.pwmControllerRight
}

func (c *Container) PWMControllerLeft() *pwm.PWM {
	if c.pwmControllerLeft == nil {
		pwmPinLeft := gpioreg.ByName(c.cfg.PWMPinLeft)
		if pwmPinLeft == nil {
			panic("failed to find pwm pin left")
		}

		c.pwmControllerLeft = pwm.NewPWM(periph.NewPin(pwmPinLeft))
	}

	return c.pwmControllerLeft
}

func (c *Container) Chassis() *chassis.Chassis {
	c.chassisMu.Lock()
	defer c.chassisMu.Unlock()

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

		c.chassis = chassis.NewChassis(
			c.PWMControllerRight(),
			c.PWMControllerLeft(),
			periph.NewPin(controlPinRight1),
			periph.NewPin(controlPinRight2),
			periph.NewPin(controlPinLeft1),
			periph.NewPin(controlPinLeft2),
			c.cfg.PWMFrequency,
			time.Duration(c.cfg.MoveDuration)*time.Millisecond,
		)
	}

	return c.chassis
}

func (c *Container) Logger() *slog.Logger {
	if c.logger == nil {
		c.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	return c.logger
}
