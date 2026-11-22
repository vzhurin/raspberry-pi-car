package joystick

import (
	"context"
	"errors"
	"log/slog"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const delay = 1000 * time.Millisecond

type Watcher struct {
	joystickFactory func() (*Joystick, error)
	logger          *slog.Logger
	joystickID      int
	wg              sync.WaitGroup
	done            chan struct{}
	errNotify       chan error
}

func NewWatcher(
	joystickFactory func() (*Joystick, error),
	logger *slog.Logger,
	joystickID int,
) *Watcher {
	return &Watcher{
		joystickFactory: joystickFactory,
		logger:          logger,
		joystickID:      joystickID,
		done:            make(chan struct{}),
		errNotify:       make(chan error, 1),
	}
}

func (w *Watcher) Watch(ctx context.Context) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		defer close(w.errNotify)
		w.errNotify <- w.watch(ctx)
	}()
}

func (w *Watcher) ErrNotify() <-chan error {
	return w.errNotify
}

func (w *Watcher) watch(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer func() { _ = watcher.Close() }()

	if err := watcher.Add("/dev/input"); err != nil {
		return err
	}

	deviceName := "js" + strconv.Itoa(w.joystickID)

	var current *Joystick
	var joystickResult <-chan error

	stopJoystick := func() {
		if current != nil {
			current.Shutdown()
			current = nil
		}

		joystickResult = nil
	}
	defer stopJoystick()

	for {
		select {
		case <-w.done:
			return nil

		case <-ctx.Done():
			return ctx.Err()

		case err, ok := <-joystickResult:
			stopJoystick()

			if !ok || err == nil {
				continue
			}

			if errors.Is(err, ErrDriver) {
				w.logger.With(slog.Any("error", err)).Warn("joystick unavailable; waiting for reconnection")
				continue
			}

			return err
		case event, ok := <-watcher.Events:
			if !ok {
				return errors.New("watcher events channel closed")
			}

			if filepath.Base(event.Name) != deviceName {
				continue
			}

			if event.Has(fsnotify.Remove) {
				stopJoystick()
			}

			if event.Has(fsnotify.Create) && current == nil {
				time.Sleep(delay)

				next, err := w.joystickFactory()
				if err != nil {
					return err
				}

				current = next
				current.Run(ctx)
				joystickResult = current.ErrNotify()
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("watcher errors channel closed")
			}

			w.logger.Error(
				"error watching directory",
				slog.Any("error", err),
			)
		}
	}
}

func (w *Watcher) Shutdown() {
	close(w.done)
	w.wg.Wait()
}
