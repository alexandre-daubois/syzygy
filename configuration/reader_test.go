package configuration

import "testing"

func TestNewReader(t *testing.T) {
	r := NewReader()
	if r == nil {
		t.Error("NewReader() should not return nil")
	}
}

func TestYamlReaderRead(t *testing.T) {
	r := NewReader()
	c, err := r.Read("../fixture/simple_config.yaml")
	if err != nil {
		t.Errorf("Read() should not return error: %v", err)
	}

	if c == nil {
		t.Error("Read() should not return nil")
	}

	if len(c.Processes) < 1 {
		t.Errorf("Read() should return 1 process, got %d", len(c.Processes))
	}

	p, ok := c.Processes["process1"]
	if !ok {
		t.Error("Read() should return process with name 'process1'")
	}

	if p.Command != "echo hello world" {
		t.Errorf("Read() should return process with command 'sleep 1', got '%s'", p.Command)
	}

	if p.Cwd != "/tmp" {
		t.Errorf("Read() should return process with cwd '/tmp', got '%s'", p.Cwd)
	}

	if len(p.Env) != 1 {
		t.Errorf("Read() should return process with 1 env var, got %d", len(p.Env))
	}

	if p.Env[0] != "FOO=bar" {
		t.Errorf("Read() should return process with env var 'FOO=bar', got '%s'", p.Env[0])
	}

	if p.StopSignal != "SIGTERM" {
		t.Errorf("Read() should return process with stop signal 'SIGTERM', got '%s'", p.StopSignal)
	}

	if p.RestartPolicy != "never" {
		t.Errorf("Read() should return process with restart policy 'always', got '%s'", p.RestartPolicy)
	}

	if p.OutputLogFile != "/tmp/process1.out.log" {
		t.Errorf("Read() should return process with output log file '/tmp/process1.out.log', got '%s'", p.OutputLogFile)
	}

	if p.EventsLogFile != "/tmp/process1.events.log" {
		t.Errorf("Read() should return process with events log file '/tmp/process1.events.log', got '%s'", p.EventsLogFile)
	}
}

func TestYamlReaderReadNotExistingFile(t *testing.T) {
	r := NewReader()
	_, err := r.Read("not_existing_file.yaml")
	if err == nil {
		t.Error("Read() should return error")
	}
}

func TestYamlReaderCannotUnmarshall(t *testing.T) {
	r := NewReader()
	_, err := r.Read("../fixture/invalid_config.yaml")
	if err == nil {
		t.Error("Read() should return error")
	}
}
