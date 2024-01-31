package process

import (
	"errors"
	"log"
	"os"
	"szg/configuration"
)

type ProcessManager struct {
	Processes  map[int]*Process
	LogsWriter *log.Logger

	events chan Event
}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		events:     make(chan Event, 1),
		Processes:  make(map[int]*Process),
		LogsWriter: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (pm *ProcessManager) SetLogsPath(path string) {
	if path != "" {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		pm.LogsWriter = log.New(file, "", log.LstdFlags)
	} else {
		panic("Logs path cannot be empty")
	}
}

func (pm *ProcessManager) RunFromConfiguration(configuration *configuration.Configuration) {
	for name, pc := range configuration.Processes {
		pm.RunCommand(name, &pc)
	}
}

func (pm *ProcessManager) StopAll() {
	for _, p := range pm.Processes {
		p.Stop()
	}
}

func (pm *ProcessManager) GetProcessByName(name string) *Process {
	for _, p := range pm.Processes {
		if p.Name == name {
			return p
		}
	}

	return nil
}

func (pm *ProcessManager) Stop(name string) error {
	for _, p := range pm.Processes {
		if p.Name == name {
			p.Stop()

			return nil
		}
	}

	return errors.New("process not found")
}

func (pm *ProcessManager) AddProcess(p *Process) {
	pm.Processes[p.Pid] = p
}

func (pm *ProcessManager) RemoveProcess(p *Process) {
	delete(pm.Processes, p.Pid)
}

func (pm *ProcessManager) RunCommand(name string, processConfiguration *configuration.ProcessConfiguration) {
	pm.runProcess(NewProcess(name, processConfiguration, pm.LogsWriter))
}

func (pm *ProcessManager) runProcess(p *Process) {
	err := p.Start()
	if err != nil {
		pm.LogsWriter.Printf("Cannot start process '%s' because of %s\n", p.Command, err)

		return
	}

	pm.AddProcess(p)

	go p.WatchState(pm.events)
}

func (pm *ProcessManager) handleProcessStopped(p *Process) {
	pm.LogsWriter.Printf("Process %d stopped", p.Pid)
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
				delete(pm.Processes, processEvent.Process.Pid)
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
