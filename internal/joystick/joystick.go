package joystick

import (
	"context"
	"errors"
)

var ErrDriver = errors.New("joystick driver error")

type Joystick interface {
	Run(context.Context)
	Shutdown()
	Bus() <-chan State
	ErrNotify() <-chan error
}

type State struct {
	RightGas float64
	LeftGas  float64
}
