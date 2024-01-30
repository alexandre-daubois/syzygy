package process

import (
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
	p := NewProcess([]string{"echo", "hello"}, Never)

	pm.AddProcess(p)

	if pm.Processes[p.Pid] == nil {
		t.Errorf("Expected %s, got nil", "Process")
	}
}

func TestProcessManager_RemoveProcess(t *testing.T) {
	pm := NewProcessManager()
	p := NewProcess([]string{"echo", "hello"}, Never)

	pm.AddProcess(p)
	pm.RemoveProcess(p)

	if pm.Processes[p.Pid] != nil {
		t.Errorf("Expected nil, got %s", "Process")
	}
}

func TestProcessManager_RunCommand(t *testing.T) {
	pm := NewProcessManager()

	pm.RunCommand([]string{"echo", "hello"}, Never)

	if len(pm.Processes) != 1 {
		t.Errorf("Expected %s, got nil", "Process")
	}
}

func TestProcessManager_Loop(t *testing.T) {
	pm := NewProcessManager()

	pm.RunCommand([]string{"echo", "hello"}, Never)
	pm.RunCommand([]string{"echo", "hello2"}, Never)

	go pm.Loop()

	time.Sleep(20 * time.Millisecond)

	if len(pm.Processes) != 0 {
		t.Errorf("Expected nil, got %s", "Process")
	}
}
