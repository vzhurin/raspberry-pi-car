package hardware

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func DeviceID(path string) (int, error) {
	base := filepath.Base(path)
	i := len(base)
	for i > 0 && base[i-1] >= '0' && base[i-1] <= '9' {
		i--
	}

	if i == len(base) {
		return -1, fmt.Errorf("path %q does not contain device ID", path)
	}

	return strconv.Atoi(base[i:])
}

func CutDevice(device string) string {
	split := strings.Split(device, "/")
	return strings.Join(split[:len(split)-1], "/")
}
