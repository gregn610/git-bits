package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/nerdalize/git-bits/command"
)

var version = "0.0.0"

func main() {
	rootCmd := &cobra.Command{
		Use:     "git-bits",
		Short:   "Git extension for large binary files",
		Version: version,
	}

	// Add subcommands
	rootCmd.AddCommand(
		command.NewScanCmd(),
		command.NewSplitCmd(),
		command.NewInstallCmd(),
		command.NewFetchCmd(),
		command.NewPullCmd(),
		command.NewPushCmd(),
		command.NewCombineCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
