package command

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewScanCmd(t *testing.T) {
	cmd := NewScanCmd()
	if cmd.Use != "scan" {
		t.Errorf("Expected Use to be 'scan', got %s", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("Expected Short description to be set")
	}
}

func TestNewSplitCmd(t *testing.T) {
	cmd := NewSplitCmd()
	if cmd.Use != "split" {
		t.Errorf("Expected Use to be 'split', got %s", cmd.Use)
	}
}

func TestNewFetchCmd(t *testing.T) {
	cmd := NewFetchCmd()
	if cmd.Use != "fetch" {
		t.Errorf("Expected Use to be 'fetch', got %s", cmd.Use)
	}
}

func TestNewPullCmd(t *testing.T) {
	cmd := NewPullCmd()
	if cmd.Use != "pull" {
		t.Errorf("Expected Use to be 'pull', got %s", cmd.Use)
	}
}

func TestNewPushCmd(t *testing.T) {
	cmd := NewPushCmd()
	if cmd.Use != "push" {
		t.Errorf("Expected Use to be 'push', got %s", cmd.Use)
	}
}

func TestNewCombineCmd(t *testing.T) {
	cmd := NewCombineCmd()
	if cmd.Use != "combine" {
		t.Errorf("Expected Use to be 'combine', got %s", cmd.Use)
	}
}

func TestAllCommandsHaveHelp(t *testing.T) {
	commands := []*cobra.Command{
		NewScanCmd(),
		NewSplitCmd(),
		NewFetchCmd(),
		NewPullCmd(),
		NewPushCmd(),
		NewCombineCmd(),
	}

	for _, cmd := range commands {
		if cmd.Short == "" {
			t.Errorf("Command %s missing Short description", cmd.Use)
		}
		if cmd.RunE == nil {
			t.Errorf("Command %s missing RunE function", cmd.Use)
		}
	}
}