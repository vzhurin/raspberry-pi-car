package joystick

import (
	"context"
	"errors"
	"log/slog"
	"raspberry-pi-car/internal/hardware"
	"sync"
)

type joystickFactory func(device string) (Joystick, error)

type Dispatcher struct {
	logger    *slog.Logger
	factory   joystickFactory
	eventBus  <-chan hardware.Event
	stateBus  chan State
	wg        sync.WaitGroup
	done      chan struct{}
	errNotify chan error

	current         Joystick
	currentStateBus <-chan State
	errors          <-chan error
}

func NewDispatcher(
	logger *slog.Logger,
	factory joystickFactory,
	eventBus <-chan hardware.Event,
) *Dispatcher {
	return &Dispatcher{
		logger:    logger,
		factory:   factory,
		eventBus:  eventBus,
		stateBus:  make(chan State),
		done:      make(chan struct{}),
		errNotify: make(chan error, 1),
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
	d.logger.Info("shutting down")
	close(d.done)
	d.wg.Wait()
}

func (d *Dispatcher) StateBus() <-chan State {
	return d.stateBus
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
		case event, ok := <-d.eventBus:
			if !ok {
				d.logger.Debug("bus channel closed")
				return nil
			}

			switch e := event.(type) {
			case *hardware.Plugged:
				d.logger.Debug("dispatching plugged event")
				if err := d.onPlugged(ctx, e); err != nil {
					return err
				}
			case *hardware.Unplugged:
				d.logger.Warn("dispatching unplugged event")
				if err := d.onUnplugged(ctx, e); err != nil {
					return err
				}
			default:
				d.logger.Warn("unknown event")
			}
		case state := <-d.currentStateBus:
			select {
			case <-ctx.Done():
				d.logger.Debug("context done")
				return ctx.Err()
			case <-d.done:
				d.logger.Debug("done channel closed")
				return nil
			case d.stateBus <- state:
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

func (d *Dispatcher) onPlugged(ctx context.Context, event *hardware.Plugged) error {
	j, err := d.factory(event.Device())
	if err != nil {
		return err
	}
	j.Run(ctx)
	d.current = j
	d.currentStateBus = j.Bus()
	d.errors = j.ErrNotify()

	return nil
}

func (d *Dispatcher) onUnplugged(_ context.Context, _ *hardware.Unplugged) error {
	if d.current != nil {
		d.current.Shutdown()
	}

	d.current = nil
	d.currentStateBus = nil
	d.errors = nil

	return nil
}
