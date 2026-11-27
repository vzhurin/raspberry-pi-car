package hardware

import (
	"context"
	"errors"
	"log/slog"
	"raspberry-pi-car/internal/filesystem"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const watcherDelay = 500 * time.Millisecond

type Watcher struct {
	logger    *slog.Logger
	device    string
	bus       chan Event
	done      chan struct{}
	errNotify chan error
	wg        sync.WaitGroup
}

func NewWatcher(
	logger *slog.Logger,
	device string,
) *Watcher {
	return &Watcher{
		logger:    logger,
		device:    device,
		bus:       make(chan Event),
		done:      make(chan struct{}),
		errNotify: make(chan error, 1),
	}
}

func (w *Watcher) Watch(ctx context.Context) {
	w.wg.Go(func() {
		defer close(w.errNotify)
		w.errNotify <- w.watch(ctx)
	})
}

func (w *Watcher) ErrNotify() <-chan error {
	return w.errNotify
}

func (w *Watcher) Bus() <-chan Event {
	return w.bus
}

func (w *Watcher) watch(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer func() { _ = watcher.Close() }()

	path := CutDevice(w.device)
	if err := watcher.Add(path); err != nil {
		return err
	}
	w.logger.Debug("start watching", slog.String("path", path))

	if filesystem.FileExists(w.device) {
		if !w.publish(ctx, NewPlugged(w.device), watcherDelay) {
			return ctx.Err()
		}
		w.logger.Debug("device already plugged", slog.String("device", w.device))
	}

	for {
		select {
		case <-w.done:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case event, ok := <-watcher.Events:
			if !ok {
				return errors.New("watcher events channel closed")
			}

			if event.Name != w.device {
				continue
			}

			if event.Has(fsnotify.Create) {
				if !w.publish(ctx, NewPlugged(w.device), watcherDelay) {
					return ctx.Err()
				}
				w.logger.Debug("device has been plugged", slog.String("device", w.device))
				continue
			}

			if event.Has(fsnotify.Remove) {
				if !w.publish(ctx, NewUnplugged(w.device), watcherDelay) {
					return ctx.Err()
				}
				w.logger.Debug("device has been unplugged", slog.String("device", w.device))
				continue
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("watcher errors channel closed")
			}

			w.logger.Error("error watching directory", slog.Any("error", err))
		}
	}
}

func (w *Watcher) publish(ctx context.Context, event Event, delay time.Duration) bool {
	if delay > 0 {
		timer := time.NewTimer(delay)
		defer timer.Stop()

		select {
		case <-w.done:
			return false
		case <-ctx.Done():
			return false
		case <-timer.C:
		}
	}

	select {
	case <-w.done:
		return false
	case <-ctx.Done():
		return false
	case w.bus <- event:
		return true
	}
}

func (w *Watcher) Shutdown() {
	w.logger.Info("shutting down")
	close(w.done)
	w.wg.Wait()
}
