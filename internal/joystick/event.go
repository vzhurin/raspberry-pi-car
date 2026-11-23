package joystick

type Event interface {
	event()
}

type Plugged struct {
	JoystickID int
}

func NewPlugged(joystickID int) *Plugged {
	return &Plugged{JoystickID: joystickID}
}

func (*Plugged) event() {}

type Unplugged struct {
	JoystickID int
}

func NewUnplugged(joystickID int) *Unplugged {
	return &Unplugged{JoystickID: joystickID}
}

func (*Unplugged) event() {}
