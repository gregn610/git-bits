package main

import (
	"os"
	"testing"
)

func TestVersion(t *testing.T) {
	if version == "" {
		t.Error("Version should not be empty")
	}
}

func TestMainFunction(t *testing.T) {
	// Test that main function exists and can be called
	// We'll test with --help to avoid actual execution
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	os.Args = []string{"git-bits", "--help"}
	
	// This would normally call os.Exit, so we can't test the actual execution
	// But we can test that the function exists and compiles
	defer func() {
		if r := recover(); r != nil {
			// Expected behavior for --help
		}
	}()
	
	// Just verify the function signature exists
	main()
}