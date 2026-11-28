package _0xcafed00d

import (
	"context"
	"errors"
	"log/slog"
	"raspberry-pi-car/internal/joystick"
	"sync"
	"time"

	joy "github.com/0xcafed00d/joystick"
)

const (
	axisMax   = 32767
	rightAxis = 4
	leftAxis  = 1
)

type Joystick struct {
	driver     joy.Joystick
	logger     *slog.Logger
	bus        chan joystick.State
	pollPeriod time.Duration
	wg         sync.WaitGroup
	done       chan struct{}
	errNotify  chan error
}

func New(
	logger *slog.Logger,
	joystickID int,
	pollPeriod time.Duration,
) (*Joystick, error) {
	if joystickID < 0 {
		return nil, errors.New("invalid joystick ID")
	}

	driver, err := joy.Open(joystickID)
	if err != nil {
		return nil, errors.Join(joystick.ErrDriver, err)
	}

	return &Joystick{
		driver:     driver,
		logger:     logger,
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
	j.driver.Close()
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
		s, err := j.driver.Read()
		if err != nil {
			return errors.Join(joystick.ErrDriver, err)
		}

		right := float64(s.AxisData[rightAxis]) / axisMax
		left := float64(s.AxisData[leftAxis]) / axisMax

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
