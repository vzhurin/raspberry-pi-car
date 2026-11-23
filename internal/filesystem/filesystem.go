package filesystem

import (
	"errors"
	"os"
)

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return err == nil && !info.IsDir()
}
