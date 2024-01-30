package process

const (
	Started = iota
	Exited
	Restarted
)

type Event struct {
	Event   int
	Process *Process
}
