package command

import (
	"strings"
	"testing"
)

func TestNewInstallCmd(t *testing.T) {
	cmd := NewInstallCmd()
	
	if cmd.Use != "install" {
		t.Errorf("Expected Use to be 'install', got %s", cmd.Use)
	}
	
	if cmd.Short == "" {
		t.Error("Expected Short description to be set")
	}
	
	// Test flags exist
	bucketFlag := cmd.Flags().Lookup("bucket")
	if bucketFlag == nil {
		t.Error("Expected bucket flag to exist")
	}
	
	remoteFlag := cmd.Flags().Lookup("remote")
	if remoteFlag == nil {
		t.Error("Expected remote flag to exist")
	}
	
	if remoteFlag.DefValue != "origin" {
		t.Errorf("Expected remote flag default to be 'origin', got %s", remoteFlag.DefValue)
	}
}

func TestAskInputValidation(t *testing.T) {
	// Test with mock input - this is a simple validation test
	// Real input testing would require more complex mocking
	
	// Just test that the function exists and has correct signature
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic due to no stdin in test environment
			if !strings.Contains(r.(error).Error(), "inappropriate ioctl") {
				t.Errorf("Unexpected panic: %v", r)
			}
		}
	}()
	
	// This will panic in test environment, which is expected
	askInput("test prompt")
}