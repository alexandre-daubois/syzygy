package process

import (
	"hypervigo/configuration"
	"strings"
	"testing"
	"time"
)

func TestNewProcessManager(t *testing.T) {
	pm := NewProcessManager()

	if pm.Processes == nil {
		t.Errorf("Expected %s, got nil", "map[int]*Process")
	}

	if pm.events == nil {
		t.Errorf("Expected %s, got nil", "chan Event")
	}
}

func TestProcessManager_AddProcess(t *testing.T) {
	pm := NewProcessManager()
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "never",
	}

	p := NewProcess("process", &c)

	pm.AddProcess(p)

	if pm.Processes[p.Pid].Name != "process" {
		t.Errorf("Expected %s, got %s", "process", pm.Processes[p.Pid].Name)
	}

	if strings.Join(pm.Processes[p.Pid].Command, " ") != "echo hello" {
		t.Errorf("Expected %s, got %s", "echo", strings.Join(pm.Processes[p.Pid].Command, " "))
	}

	if pm.Processes[p.Pid] == nil {
		t.Errorf("Expected %s, got nil", "Process")
	}
}

func TestProcessManager_RemoveProcess(t *testing.T) {
	pm := NewProcessManager()
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "never",
	}

	p := NewProcess("process", &c)

	pm.AddProcess(p)
	pm.RemoveProcess(p)

	if pm.Processes[p.Pid] != nil {
		t.Errorf("Expected nil, got %s", "Process")
	}
}

func TestProcessManager_RunCommand(t *testing.T) {
	pm := NewProcessManager()
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "never",
	}

	pm.RunCommand("process", &c)

	if len(pm.Processes) != 1 {
		t.Errorf("Expected %s, got nil", "Process")
	}
}

func TestProcessManager_Loop(t *testing.T) {
	pm := NewProcessManager()
	c1 := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "never",
	}

	c2 := configuration.ProcessConfiguration{
		Command:       "echo hello2",
		RestartPolicy: "never",
	}

	pm.RunCommand("process", &c1)
	pm.RunCommand("process", &c2)

	go pm.Loop()

	time.Sleep(20 * time.Millisecond)

	if len(pm.Processes) != 0 {
		t.Errorf("Expected nil, got %s", "Process")
	}
}
