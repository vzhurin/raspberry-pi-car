package joystick

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/0xcafed00d/joystick"
)

const (
	axisMax   = 32767
	rightAxis = 4
	leftAxis  = 1
)

var ErrDriver = errors.New("joystick driver error")

type mover interface {
	Move(ctx context.Context, rightDutyCycle, leftDutyCycle float64) error
}

type Joystick struct {
	driver     joystick.Joystick
	mover      mover
	logger     *slog.Logger
	pollPeriod time.Duration
	wg         sync.WaitGroup
	done       chan struct{}
	errNotify  chan error
}

func New(
	m mover,
	logger *slog.Logger,
	joystickID int,
	pollPeriod time.Duration,
) (*Joystick, error) {
	driver, err := joystick.Open(joystickID)
	if err != nil {
		return nil, errors.Join(ErrDriver, err)
	}

	return &Joystick{
		driver:     driver,
		mover:      m,
		logger:     logger,
		pollPeriod: pollPeriod,
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

func (j *Joystick) ErrNotify() <-chan error {
	return j.errNotify
}

func (j *Joystick) run(ctx context.Context) error {
	ticker := time.NewTicker(j.pollPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			j.logger.Debug("context done")
			return ctx.Err()
		case <-j.done:
			j.logger.Debug("done channel closed")
			return nil
		case <-ticker.C:
		}

		state, err := j.driver.Read()
		if err != nil {
			return errors.Join(ErrDriver, err)
		}

		right := -float64(state.AxisData[rightAxis]) / axisMax
		left := -float64(state.AxisData[leftAxis]) / axisMax

		// right and left are mixed up on purpose
		err = j.mover.Move(ctx, left, right)
		if err != nil {
			return err
		}
	}
}
