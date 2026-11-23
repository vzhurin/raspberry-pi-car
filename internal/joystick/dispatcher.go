package joystick

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

type joystickFactory func(id int) (*Joystick, error)

type Dispatcher struct {
	logger    *slog.Logger
	factory   joystickFactory
	bus       <-chan Event
	wg        sync.WaitGroup
	done      chan struct{}
	errNotify chan error

	current *Joystick
	errors  <-chan error
}

func NewDispatcher(logger *slog.Logger, factory joystickFactory, bus <-chan Event) *Dispatcher {
	return &Dispatcher{
		logger:    logger,
		factory:   factory,
		bus:       bus,
		done:      make(chan struct{}),
		errNotify: make(chan error),
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context) {
	d.wg.Go(func() {
		defer close(d.errNotify)
		d.errNotify <- d.dispatch(ctx)
	})
}

func (d *Dispatcher) ErrNotify() <-chan error {
	return d.errNotify
}

func (d *Dispatcher) Shutdown() {
	d.logger.Info("shutting down dispatcher")
	close(d.done)
	d.wg.Wait()
}

func (d *Dispatcher) dispatch(ctx context.Context) error {
	stopJoystick := func() {
		if d.current != nil {
			d.logger.Debug("shutting down joystick")
			d.current.Shutdown()
			d.current = nil
		}

		d.errors = nil
	}
	defer stopJoystick()

	for {
		select {
		case <-ctx.Done():
			d.logger.Debug("context done")
			return ctx.Err()
		case <-d.done:
			d.logger.Debug("done channel closed")
			return nil
		case event, ok := <-d.bus:
			if !ok {
				d.logger.Debug("bus channel closed")
				return nil
			}

			switch e := event.(type) {
			case *Plugged:
				d.logger.Debug("dispatching plugged event")
				if err := d.onPlugged(ctx, e); err != nil {
					d.errNotify <- err
				}
			case *Unplugged:
				if err := d.onUnplugged(e); err != nil {
					d.logger.Warn("dispatching unplugged event")
					d.errNotify <- err
				}
			default:
				d.logger.Warn("unknown event")
			}
		case err := <-d.errors:
			if errors.Is(err, ErrDriver) {
				d.logger.Warn("joystick driver error", slog.Any("error", err))
				stopJoystick()
				continue
			}

			return err
		}
	}
}

func (d *Dispatcher) onPlugged(ctx context.Context, event *Plugged) error {
	j, err := d.factory(event.JoystickID)
	if err != nil {
		return err
	}
	j.Run(ctx)
	d.current = j
	d.errors = j.errNotify

	return nil
}

func (d *Dispatcher) onUnplugged(_ *Unplugged) error {
	if d.current != nil {
		d.current.Shutdown()
	}

	d.current = nil
	d.errors = nil

	return nil
}
