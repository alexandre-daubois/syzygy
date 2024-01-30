package process

import (
	"log"
	"os"
	"os/exec"
)

type Process struct {
	Command       []string
	Pid           int
	RestartPolicy int
	StartCount    int

	process *os.Process
}

func NewProcess(command []string, restartPolicy int) *Process {
	return &Process{
		Command:       command,
		RestartPolicy: restartPolicy,
		StartCount:    0,
	}
}

func (p *Process) Start() error {
	cmd := exec.Command(p.Command[0], p.Command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatalf("Cannot start process '%s' because of %s\n", p.Command, err)
	}

	p.process = cmd.Process
	p.Pid = cmd.Process.Pid
	p.StartCount++

	log.Printf("Process '%s' started with pid %d\n", p.Command, p.Pid)

	return nil
}

func (p *Process) Stop() error {
	go func() {
		err := p.process.Signal(os.Interrupt)

		if err != nil {
			log.Fatalf("Cannot stop process '%s' because of %s\n", p.Command, err)

			return
		}

		p.process.Wait()
	}()

	return nil
}

func (p *Process) WatchState(events chan Event) {
	state, err := p.process.Wait()
	if err != nil {
		log.Fatalf("Cannot watch state of process '%s' because of %s\n", p.Command, err)
	}

	if state.Exited() {
		p.handleExit(events)
	}
}

func (p *Process) handleExit(events chan Event) {
	switch p.RestartPolicy {
	case Always:
		log.Printf("restarting '%s' (always)\n", p.Command)
		events <- Event{Event: Restarted, Process: p}
	default:
		events <- Event{Event: Exited, Process: p}
		log.Printf("'%s' is not configured to restart\n", p.Command)
	}
}
