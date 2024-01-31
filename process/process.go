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

	process      *os.Process
	eventsWriter *log.Logger
}

func NewProcess(name string, processConfiguration *configuration.ProcessConfiguration) *Process {
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

	return &Process{
		Name:          strings.TrimSpace(name),
		Command:       strings.Split(processConfiguration.Command, " "),
		RestartPolicy: restartPolicy,
		StartCount:    0,
		Configuration: processConfiguration,
	}
}

func (p *Process) Start() error {
	cmd := exec.Command(p.Command[0], p.Command[1:]...)

	var output, errout *os.File
	var err error

	if p.Configuration.EventsLogFile != "" {
		file, err := os.OpenFile(p.Configuration.EventsLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Cannot create file '%s' because of %s\n", p.Configuration.EventsLogFile, err)

			return err
		}

		p.eventsWriter = log.New(file, "", log.LstdFlags)
	} else {
		p.eventsWriter = log.New(os.Stdout, "", log.LstdFlags)
	}

	if p.Configuration.OutputLogFile != "" {
		output, err = os.OpenFile(p.Configuration.OutputLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			p.eventsWriter.Printf("Cannot create file '%s' because of %s\n", p.Configuration.OutputLogFile, err)

			return err
		}

		errout = output
	} else {
		output = os.Stdout
		errout = os.Stderr
	}

	cmd.Stdout = output
	cmd.Stderr = errout

	if p.Configuration.Cwd != "" {
		cmd.Dir = p.Configuration.Cwd
	}

	if len(p.Configuration.Env) > 0 {
		cmd.Env = append(os.Environ(), p.Configuration.Env...)
	}

	err = cmd.Start()
	if err != nil {
		p.eventsWriter.Printf("Cannot start process '%s' because of %s\n", p.Command, err)

		return err
	}

	p.process = cmd.Process
	p.Pid = cmd.Process.Pid
	p.StartCount++

	p.eventsWriter.Printf("Process '%s' started with pid %d\n", p.Command, p.Pid)

	return nil
}

// Stop is not implemented nor tested yet
func (p *Process) Stop() error {
	go func() {
		err := p.process.Signal(os.Interrupt)

		if err != nil {
			p.eventsWriter.Fatalf("Cannot stop process '%s' because of %s\n", p.Command, err)

			return
		}

		p.process.Wait()
	}()

	return nil
}

func (p *Process) WatchState(events chan Event) {
	state, err := p.process.Wait()
	if err != nil {
		p.eventsWriter.Fatalf("Cannot watch state of process '%s' because of %s\n", p.Command, err)
	}

	if state.Exited() {
		p.handleExit(events)
	}
}

func (p *Process) handleExit(events chan Event) {
	switch p.RestartPolicy {
	case Always:
		p.eventsWriter.Printf("restarting '%s' (always)\n", p.Command)
		events <- Event{Event: Restarted, Process: p}
	default:
		events <- Event{Event: Exited, Process: p}
		p.eventsWriter.Printf("'%s' is not configured to restart\n", p.Command)
	}
}
