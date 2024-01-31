package process

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"szg/configuration"
)

type Process struct {
	Name          string
	Command       []string
	Pid           int
	RestartPolicy int
	StartCount    int
	Configuration *configuration.ProcessConfiguration
	Logger        *log.Logger
	LogOutput     bool

	process *os.Process
}

func NewProcess(name string, processConfiguration *configuration.ProcessConfiguration, logger *log.Logger) *Process {
	var restartPolicy int
	switch processConfiguration.RestartPolicy {
	case "always":
		restartPolicy = Always
	case "never":
		restartPolicy = Never
	case "":
		// be default, never restart the process
		restartPolicy = Never
	default:
		panic("Unknown restart policy for process " + name)
	}

	if strings.TrimSpace(name) == "" {
		panic("Process name cannot be empty")
	}

	if processConfiguration.StopSignal != "" &&
		processConfiguration.StopSignal != "SIGINT" &&
		processConfiguration.StopSignal != "SIGKILL" {
		panic("Process stop signal must be either SIGINT or SIGKILL")
	}

	return &Process{
		Name:          strings.TrimSpace(name),
		Command:       strings.Split(processConfiguration.Command, " "),
		RestartPolicy: restartPolicy,
		StartCount:    0,
		Configuration: processConfiguration,
		Logger:        logger,
	}
}

func (p *Process) Start() error {
	cmd := exec.Command(p.Command[0], p.Command[1:]...)

	if p.Configuration.Cwd != "" {
		cmd.Dir = p.Configuration.Cwd
	}

	if len(p.Configuration.Env) > 0 {
		cmd.Env = append(os.Environ(), p.Configuration.Env...)
	}

	if p.LogOutput {
		cmd.Stdout = p.Logger.Writer()
		cmd.Stderr = p.Logger.Writer()
	}

	err := cmd.Start()
	if err != nil {
		p.Logger.Printf("Cannot start process '%s' because of %s\n", p.Command, err)

		return err
	}

	p.process = cmd.Process
	p.Pid = cmd.Process.Pid
	p.StartCount++

	p.Logger.Printf("Process '%s' started with pid %d\n", p.Command, p.Pid)

	return nil
}

func (p *Process) Stop() error {
	go func() {
		stopSignal := p.Configuration.StopSignal
		var signalToSend os.Signal

		switch stopSignal {
		case "SIGINT":
			signalToSend = os.Interrupt
		case "SIGKILL":
			signalToSend = os.Kill
		default:
			signalToSend = os.Interrupt
		}

		err := p.process.Signal(signalToSend)

		if err != nil {
			p.Logger.Fatalf("Cannot stop process '%s' because of %s\n", p.Command, err)

			return
		}

		p.process.Wait()
		p.Logger.Printf("Process '%s' stopped\n", p.Command)
	}()

	return nil
}

func (p *Process) WatchState(events chan Event) {
	state, err := p.process.Wait()
	if err != nil {
		p.Logger.Fatalf("Cannot watch state of process '%s' because of %s\n", p.Command, err)
	}

	if state.Exited() {
		p.handleExit(events)
	}
}

func (p *Process) handleExit(events chan Event) {
	switch p.RestartPolicy {
	case Always:
		p.Logger.Printf("restarting '%s' (always)\n", p.Command)
		events <- Event{Event: Restarted, Process: p}
	default:
		events <- Event{Event: Exited, Process: p}
		p.Logger.Printf("'%s' is not configured to restart\n", p.Command)
	}
}
