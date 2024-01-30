package process

import (
	"fmt"
	"strconv"
	"time"
)

type ProcessManager struct {
	Processes map[int]*Process

	events chan ProcessEvent
}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		events:    make(chan ProcessEvent, 1),
		Processes: make(map[int]*Process),
	}
}

func (pm *ProcessManager) AddProcess(p *Process) {
	pm.Processes[p.Pid] = p
}

func (pm *ProcessManager) RunCommand(command ...string) {
	p := NewProcess(pm.events, command...)
	pm.AddProcess(p)

	p.Start()
}

// Watch checks for process events to happen and notify the event channel
func (pm *ProcessManager) Watch() {
	for {
		time.Sleep(1 * time.Second)

		for _, p := range pm.Processes {
			if p.State().Exited() {
				delete(pm.Processes, p.Pid)
				p.events <- ProcessEvent{Event: Stopped, Pid: p.Pid}
			}
		}
	}
}

func (pm *ProcessManager) Loop() {
	defer close(pm.events)
	go pm.Watch()

	for {
		select {
		case processEvent := <-pm.events:
			switch processEvent.Event {
			case Started:
				fmt.Println("Process " + strconv.Itoa(processEvent.Pid) + " started")
			case Stopped:
				fmt.Println("Process " + strconv.Itoa(processEvent.Pid) + " stopped")
			}
		}

		if len(pm.Processes) == 0 {
			break
		}
	}
}
