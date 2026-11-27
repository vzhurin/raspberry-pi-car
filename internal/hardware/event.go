package hardware

type Event interface {
	Device() string
}

type Plugged struct {
	device string
}

func NewPlugged(device string) *Plugged {
	return &Plugged{device: device}
}

func (e *Plugged) Device() string {
	return e.device
}

type Unplugged struct {
	device string
}

func NewUnplugged(device string) *Unplugged {
	return &Unplugged{device: device}
}

func (e *Unplugged) Device() string {
	return e.device
}
