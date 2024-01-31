package main

import (
	"os"
	"testing"
)

func TestMainWithoutOsArgs(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, got nil")
		}
	}()

	main()
}

func TestMainWithNonExistingConfigurationFile(t *testing.T) {
	os.Args = []string{"", "non-existing-file"}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, got nil")
		}
	}()

	main()
}
