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
}

func NewProcess(command ...string) *Process {
	return &Process{
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

	log.Printf("Process '%s' started with pid %d\n", p.Command, p.Pid)

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
	}()

	return nil
}

func (p *Process) WatchState(events chan Event) {
	state, err := p.process.Wait()
	if err != nil {
		log.Fatalf("Process '%s' state failed with %s\n", p.Command, err)
	}

	if state.Exited() {
		events <- Event{Event: Stopped, Pid: p.Pid}
	}
}
