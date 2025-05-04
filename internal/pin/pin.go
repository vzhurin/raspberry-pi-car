package pin

type Level bool

const (
	High Level = true
	Low  Level = false
)

type Pin interface {
	Out(level Level) error
}
