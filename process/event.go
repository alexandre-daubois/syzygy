package process

const (
	Started = iota
	Stopped
)

type Event struct {
	Event int
	Pid   int
}
