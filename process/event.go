package process

const (
	Started = iota
	Exited
	Restarted
	Stopped
)

type Event struct {
	Event   int
	Process *Process
}
