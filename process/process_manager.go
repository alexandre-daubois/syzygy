package process

import (
	"log"
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

func (pm *ProcessManager) RemoveProcess(p *Process) {
	delete(pm.Processes, p.Pid)
}

func (pm *ProcessManager) RunCommand(command ...string) {
	p := NewProcess(command...)
	p.Start()

	pm.AddProcess(p)

	go p.WatchState(pm.events)
}

func (pm *ProcessManager) handleProcessStopped(p *Process) {
	pm.RemoveProcess(p)
	log.Printf("Process %d stopped", p.Pid)
}

func (pm *ProcessManager) Loop() {
	defer close(pm.events)

	for {
		select {
		case processEvent := <-pm.events:
			switch processEvent.Event {
			case Stopped:
				pm.handleProcessStopped(pm.Processes[processEvent.Pid])
			}
		}

		if len(pm.Processes) == 0 {
			break
		}
	}
}
