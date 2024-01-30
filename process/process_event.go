package process

const (
	Started = iota
	Stopped
)

type ProcessEvent struct {
	Event int
	Pid   int
}
