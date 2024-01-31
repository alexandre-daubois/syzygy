package process

import (
	"os"
	"strings"
	"szg/configuration"
	"testing"
)

func TestProcessStart(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		RestartPolicy: "always",
	}
	p := NewProcess("process", &c)

	err := p.Start()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if p.Pid == 0 {
		t.Errorf("Expected %d, got %d", 0, p.Pid)
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
	p := NewProcess("process", &c)

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

	NewProcess("process", &c)
}

func TestNewProcessWithoutRestartPolicy(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command: "echo hello",
	}

	p := NewProcess("process", &c)

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

	NewProcess("", &c)
}

func TestProcessStartWithCwd(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "pwd",
		Cwd:           "/tmp",
		RestartPolicy: "always",
		OutputLogFile: "/tmp/process.out.log",
	}
	p := NewProcess("process", &c)

	err := p.Start()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	p.process.Wait()
	data, _ := os.ReadFile("/tmp/process.out.log")
	if string(data) != "/tmp\n" {
		t.Errorf("Expected %s, got %s", "/tmp\n", string(data))
	}

	os.Remove("/tmp/process.out.log")
}

func TestProcessStartWithEnvVars(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "env",
		Env:           []string{"FOO=bar"},
		RestartPolicy: "always",
		OutputLogFile: "/tmp/process.out.log",
	}
	p := NewProcess("process", &c)

	err := p.Start()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	p.process.Wait()
	data, _ := os.ReadFile("/tmp/process.out.log")
	if !strings.Contains(string(data), "FOO=bar\n") {
		t.Errorf("Expected %s, got %s", "FOO=bar\n", string(data))
	}

	os.Remove("/tmp/process.out.log")
}

func TestProcessStartWithOutputLogFile(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		OutputLogFile: "/tmp/process.out.log",
		RestartPolicy: "always",
	}
	p := NewProcess("process", &c)

	err := p.Start()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	p.process.Wait()
	data, _ := os.ReadFile("/tmp/process.out.log")
	if string(data) != "hello\n" {
		t.Errorf("Expected %s, got %s", "hello\n", string(data))
	}

	os.Remove("/tmp/process.out.log")
}

func TestProcessStartWithEventsLogFile(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		EventsLogFile: "/tmp/process.events.log",
		RestartPolicy: "always",
	}
	p := NewProcess("process", &c)

	err := p.Start()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	p.process.Wait()
	data, _ := os.ReadFile("/tmp/process.events.log")
	if !strings.Contains(string(data), "'[echo hello]' started with pid") {
		t.Errorf("Expected %s, got %s", "process started", string(data))
	}

	os.Remove("/tmp/process.events.log")
}

func TestProcessStartWithInvalidEventsLogFile(t *testing.T) {
	c := configuration.ProcessConfiguration{
		Command:       "echo hello",
		EventsLogFile: "/invalid/path/process.events.log",
		RestartPolicy: "always",
	}
	p := NewProcess("process", &c)

	err := p.Start()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
