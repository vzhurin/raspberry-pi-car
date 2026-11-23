package joystick

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"raspberry-pi-car/internal/filesystem"
	"strconv"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const watcherDelay = 500 * time.Millisecond

type Watcher struct {
	logger     *slog.Logger
	bus        chan Event
	joystickID int
	wg         sync.WaitGroup
	done       chan struct{}
	errNotify  chan error
}

func NewWatcher(
	logger *slog.Logger,
	bus chan Event,
	joystickID int,
) *Watcher {
	return &Watcher{
		logger:     logger,
		bus:        bus,
		joystickID: joystickID,
		done:       make(chan struct{}),
		errNotify:  make(chan error, 1),
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
	path := "/dev/input"
	deviceName := "js" + strconv.Itoa(w.joystickID)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer func() { _ = watcher.Close() }()

	if err := watcher.Add(path); err != nil {
		return err
	}

	if filesystem.FileExists(fmt.Sprintf("%s/%s", path, deviceName)) {
		w.bus <- NewPlugged(w.joystickID)
		w.logger.Debug("joystick already plugged", slog.Int("ID", w.joystickID))
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

			if filepath.Base(event.Name) != deviceName {
				continue
			}

			if event.Has(fsnotify.Create) {
				time.Sleep(watcherDelay)
				w.bus <- NewPlugged(w.joystickID)
				w.logger.Debug("joystick has been plugged", slog.Int("ID", w.joystickID))
				continue
			}

			if event.Has(fsnotify.Remove) {
				time.Sleep(watcherDelay)
				w.bus <- NewUnplugged(w.joystickID)
				w.logger.Debug("joystick has been unplugged", slog.Int("ID", w.joystickID))
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

func (w *Watcher) Shutdown() {
	w.logger.Info("shutting down watcher")
	close(w.done)
	w.wg.Wait()
}
