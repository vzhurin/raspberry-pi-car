package goSDL3

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"raspberry-pi-car/internal/joystick"
	"sync"
	"time"

	"github.com/Zyko0/go-sdl3/sdl"
)

const axisMax = 32767

type Joystick struct {
	gamepad    *sdl.Gamepad
	logger     *slog.Logger
	bus        chan joystick.State
	pollPeriod time.Duration
	wg         sync.WaitGroup
	done       chan struct{}
	errNotify  chan error
}

func New(
	logger *slog.Logger,
	gamepadID int,
	pollPeriod time.Duration,
) (*Joystick, error) {
	if gamepadID < 0 {
		return nil, errors.New("invalid gamepad ID")
	}

	gamepadsIDs, err := sdl.GetGamepads()
	if err != nil {
		return nil, err
	}

	if len(gamepadsIDs) <= gamepadID {
		return nil, fmt.Errorf("gamepad %d not found", gamepadID)
	}

	gamepad, err := gamepadsIDs[gamepadID].OpenGamepad()
	if err != nil {
		return nil, errors.Join(joystick.ErrDriver, err)
	}

	return &Joystick{
		logger:     logger,
		gamepad:    gamepad,
		pollPeriod: pollPeriod,
		bus:        make(chan joystick.State),
		done:       make(chan struct{}),
		errNotify:  make(chan error, 1),
	}, nil
}

func (j *Joystick) Run(ctx context.Context) {
	j.wg.Go(func() {
		defer close(j.errNotify)
		j.errNotify <- j.run(ctx)
	})
}

func (j *Joystick) Shutdown() {
	j.logger.Info("shutting down")
	close(j.done)
	j.gamepad.Close()
	j.wg.Wait()
}

func (j *Joystick) Bus() <-chan joystick.State {
	return j.bus
}

func (j *Joystick) ErrNotify() <-chan error {
	return j.errNotify
}

func (j *Joystick) run(ctx context.Context) error {
	ticker := time.NewTicker(j.pollPeriod)
	defer ticker.Stop()

	for {
		sdl.UpdateGamepads()

		fmt.Println(j.gamepad.Axis(sdl.GAMEPAD_AXIS_RIGHTY))

		right := float64(j.gamepad.Axis(sdl.GAMEPAD_AXIS_RIGHTY)) / axisMax
		left := float64(j.gamepad.Axis(sdl.GAMEPAD_AXIS_LEFTY)) / axisMax

		select {
		case j.bus <- joystick.State{
			RightGas: right,
			LeftGas:  left,
		}:
		case <-ctx.Done():
			j.logger.Debug("context done")
			return ctx.Err()
		case <-j.done:
			j.logger.Debug("done channel closed")
			return nil
		case <-ticker.C:
		}
	}
}
