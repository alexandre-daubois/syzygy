package process

import (
	"strings"
	"szg/configuration"
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

func TestProcessManager_SetLogsPath(t *testing.T) {
	pm := NewProcessManager()

	pm.SetLogsPath("/tmp/logs")

	if pm.LogsWriter == nil {
		t.Errorf("Expected %s, got nil", "LogsWriter")
	}
}

func TestProcessManager_GetProcessByName(t *testing.T) {
	pm := NewProcessManager()
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "never",
	}

	p := NewProcess("process", &c, pm.LogsWriter)

	pm.AddProcess(p)

	if pm.GetProcessByName("process") == nil {
		t.Errorf("Expected %s, got nil", "Process")
	}
}

func TestProcessManager_AddProcess(t *testing.T) {
	pm := NewProcessManager()
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "never",
	}

	p := NewProcess("process", &c, pm.LogsWriter)

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

	p := NewProcess("process", &c, pm.LogsWriter)

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

	pm.RunCommand("process1", &c1)
	pm.RunCommand("process2", &c2)

	go pm.Loop()

	time.Sleep(20 * time.Millisecond)

	if pm.GetProcessByName("process1").Status != ProcessExited {
		t.Errorf("Expected %s, got %s", ProcessExited, "Process")
	}
}
