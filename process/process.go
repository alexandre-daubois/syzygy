package process

import (
	"log"
	"os"
	"os/exec"
)

type Process struct {
	Command []string
	Pid     int

	process *os.Process
	events  chan ProcessEvent
}

func NewProcess(events chan ProcessEvent, command ...string) *Process {
	return &Process{
		events:  events,
		Command: command,
	}
}

func (p *Process) Start() error {
	cmd := exec.Command(p.Command[0], p.Command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatalf("Process '%s' start failed with %s\n", p.Command, err)
	}

	p.process = cmd.Process
	p.Pid = cmd.Process.Pid

	p.events <- ProcessEvent{Event: Started, Pid: p.Pid}

	return nil
}

func (p *Process) Stop() error {
	go func() {
		err := p.process.Signal(os.Interrupt)

		if err != nil {
			log.Fatalf("Process '%s' stop failed with %s\n", p.Command, err)

			return
		}

		p.process.Wait()

		p.events <- ProcessEvent{Event: Stopped, Pid: p.Pid}
	}()

	return nil
}

func (p *Process) State() *os.ProcessState {
	state, err := p.process.Wait()

	if err != nil {
		log.Fatalf("Process '%s' state failed with %s\n", p.Command, err)
	}

	return state
}
