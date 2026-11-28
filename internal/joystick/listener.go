package joystick

import (
	"context"
	"log/slog"
	"sync"
)

type mover interface {
	Move(ctx context.Context, rightDutyCycle, leftDutyCycle float64) error
}

type Listener struct {
	mover     mover
	logger    *slog.Logger
	bus       <-chan State
	wg        sync.WaitGroup
	done      chan struct{}
	errNotify chan error
}

func NewListener(
	mover mover,
	logger *slog.Logger,
	bus <-chan State,
) *Listener {
	return &Listener{
		mover:     mover,
		logger:    logger,
		bus:       bus,
		done:      make(chan struct{}),
		errNotify: make(chan error, 1),
	}
}

func (l *Listener) Run(ctx context.Context) {
	l.wg.Go(func() {
		defer close(l.errNotify)
		l.errNotify <- l.run(ctx)
	})
}

func (l *Listener) Shutdown() {
	l.logger.Info("shutting down")
	close(l.done)
	l.wg.Wait()
}

func (l *Listener) ErrNotify() <-chan error {
	return l.errNotify
}

func (l *Listener) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			l.logger.Debug("context done")
			return ctx.Err()
		case <-l.done:
			l.logger.Debug("done channel closed")
			return nil
		case s := <-l.bus:
			err := l.mover.Move(ctx, -s.LeftGas, -s.RightGas)
			if err != nil {
				return err
			}
		}
	}
}
