package process

import "testing"

func TestNewProcess(t *testing.T) {
	p := NewProcess([]string{"echo", "hello"}, Always)

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

func TestProcessStart(t *testing.T) {
	p := NewProcess([]string{"echo", "hello"}, Always)

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
