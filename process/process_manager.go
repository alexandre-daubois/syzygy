package process

import (
	"hypervigo/configuration"
	"log"
)

type ProcessManager struct {
	Processes map[int]*Process

	events chan Event
}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		events:    make(chan Event, 1),
		Processes: make(map[int]*Process),
	}
}

func (pm *ProcessManager) RunFromConfiguration(configuration *configuration.Configuration) {
	for name, pc := range configuration.Processes {
		pm.RunCommand(name, &pc)
	}
}

func (pm *ProcessManager) AddProcess(p *Process) {
	pm.Processes[p.Pid] = p
}

func (pm *ProcessManager) RemoveProcess(p *Process) {
	delete(pm.Processes, p.Pid)
}

func (pm *ProcessManager) RunCommand(name string, processConfiguration *configuration.ProcessConfiguration) {
	pm.runProcess(NewProcess(name, processConfiguration))
}

func (pm *ProcessManager) runProcess(p *Process) {
	err := p.Start()
	if err != nil {
		log.Printf("Cannot start process '%s' because of %s\n", p.Command, err)

		return
	}

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
			case Exited:
				pm.handleProcessStopped(pm.Processes[processEvent.Process.Pid])
			case Restarted:
				pm.runProcess(processEvent.Process)
			default:
				panic("unhandled process event type")
			}
		}

		if len(pm.Processes) == 0 {
			break
		}
	}
}
