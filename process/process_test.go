package process

import (
	"log"
	"os"
	"strings"
	"szg/configuration"
	"testing"
	"time"
)

func TestProcessStart(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "always",
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	p := NewProcess("process", &c, logger)

	err := p.Start()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if p.Pid == 0 {
		t.Errorf("Expected %d, got %d", 0, p.Pid)
	}

	if p.Status != ProcessRunning {
		t.Errorf("Expected %s, got %s", ProcessRunning, p.Status)
	}

	if p.StartCount != 1 {
		t.Errorf("Expected %d, got %d", 1, p.StartCount)
	}
}

func TestNewProcess(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "always",
	}
	p := NewProcess("process", &c, nil)

	if p.Name != "process" {
		t.Errorf("Expected %s, got %s", "process", p.Name)
	}

	if p.Configuration == nil {
		t.Errorf("Expected configuration to be not nil")
	}

	if p.Command[0] != "echo" {
		t.Errorf("Expected %s, got %s", "echo", p.Command[0])
	}

	if p.Command[1] != "hello" {
		t.Errorf("Expected %s, got %s", "hello", p.Command[1])
	}

	if p.RestartPolicy != Always {
		t.Errorf("Expected %d, got %d", 0, p.RestartPolicy)
	}

	if p.Status != ProcessPending {
		t.Errorf("Expected %s, got %s", ProcessPending, p.Status)
	}

	if p.StartCount != 0 {
		t.Errorf("Expected %d, got %d", 0, p.StartCount)
	}
}

func TestNewProcessWithInvalidRestartPolicy(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "invalid",
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, got nil")
		}
	}()

	NewProcess("process", &c, nil)
}

func TestNewProcessWithoutRestartPolicy(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command: "echo hello",
	}

	p := NewProcess("process", &c, nil)

	if p.RestartPolicy != Never {
		t.Errorf("Expected %d, got %d", Never, p.RestartPolicy)
	}
}

func TestNewProcessWithoutName(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "always",
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, got nil")
		}
	}()

	NewProcess("", &c, nil)
}

func TestNewProcessWithInvalidStopSignal(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "always",
		StopSignal:    "SIGTERM",
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, got nil")
		}
	}()

	NewProcess("process", &c, nil)
}

func TestProcess_Stop(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "never",
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	p := NewProcess("process", &c, logger)

	err := p.Start()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	p.Stop()

	time.Sleep(20 * time.Millisecond)

	if p.Status != ProcessStopped {
		t.Errorf("Expected %s, got %s", ProcessStopped, p.Status)
	}
}

func TestProcessStartWithCwd(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "pwd",
		Cwd:           "/tmp",
		RestartPolicy: "always",
	}

	file, _ := os.Create("/tmp/process.out.log")
	logger := log.New(file, "", log.LstdFlags)
	p := NewProcess("process", &c, logger)

	err := p.Start()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	p.process.Wait()
	data, _ := os.ReadFile("/tmp/process.out.log")
	if strings.Contains(string(data), "/tmp\n") {
		t.Errorf("Expected %s, got %s", "/tmp\n", string(data))
	}

	os.Remove("/tmp/process.out.log")
}

func TestProcessStartWithEnvVars(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "env",
		Env:           []string{"FOO=bar"},
		RestartPolicy: "always",
	}

	file, _ := os.Create("/tmp/process.out.log")
	logger := log.New(file, "", log.LstdFlags)
	p := NewProcess("process", &c, logger)

	p.LogOutput = true

	err := p.Start()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	p.process.Wait()
	data, _ := os.ReadFile("/tmp/process.out.log")
	if !strings.Contains(string(data), "FOO=bar") {
		t.Errorf("Expected %s, got %s", "FOO=bar\n", string(data))
	}

	os.Remove("/tmp/process.out.log")
}

func TestProcess_HandleExit_Always(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "always",
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	p := NewProcess("process", &c, logger)

	err := p.Start()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	events := make(chan Event, 1)
	p.HandleExit(events)

	event := <-events

	if event.Process.Pid != p.Pid {
		t.Errorf("Expected %d, got %d", p.Pid, p.Pid)
	}

	if event.Event != Restarted {
		t.Errorf("Expected %d, got %d", Exited, event.Event)
	}
}

func TestProcess_HandleExit_Never(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "never",
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	p := NewProcess("process", &c, logger)

	err := p.Start()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	events := make(chan Event, 1)
	p.HandleExit(events)

	event := <-events

	if event.Process.Pid != p.Pid {
		t.Errorf("Expected %d, got %d", p.Pid, p.Pid)
	}

	if event.Event != Exited {
		t.Errorf("Expected %d, got %d", Exited, event.Event)
	}
}

func TestProcess_HandleExit_UnlessStopped_Exited(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "unless-stopped",
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	p := NewProcess("process", &c, logger)

	err := p.Start()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	events := make(chan Event, 1)
	p.HandleExit(events)

	event := <-events

	if event.Process.Pid != p.Pid {
		t.Errorf("Expected %d, got %d", p.Pid, p.Pid)
	}

	if event.Event != Restarted {
		t.Errorf("Expected %d, got %d", Restarted, event.Event)
	}
}

func TestProcess_HandleExit_UnlessStopped_Stopped(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "unless-stopped",
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	p := NewProcess("process", &c, logger)

	err := p.Start()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	p.Status = ProcessStopped

	events := make(chan Event, 1)
	p.HandleExit(events)

	event := <-events

	if event.Process.Pid != p.Pid {
		t.Errorf("Expected %d, got %d", p.Pid, p.Pid)
	}

	if event.Event != Stopped {
		t.Errorf("Expected %d, got %d", Stopped, event.Event)
	}

	close(events)
}
