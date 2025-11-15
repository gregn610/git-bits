package command

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/nerdalize/git-bits/bits"
)

func NewScanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scan",
		Short: "queries the git database for all chunk keys in blobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, _ := os.Getwd()
			repo, err := bits.NewRepository(wd, os.Stderr)
			if err != nil {
				return err
			}
			return repo.ScanEach(os.Stdin, os.Stdout)
		},
	}
}

func NewSplitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "split",
		Short: "splits a file into chunks and store them locally",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, _ := os.Getwd()
			repo, err := bits.NewRepository(wd, os.Stderr)
			if err != nil {
				return err
			}
			return repo.Split(os.Stdin, os.Stdout)
		},
	}
}

func NewFetchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fetch",
		Short: "fetch chunks from the remote store and save each locally",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, _ := os.Getwd()
			repo, err := bits.NewRepository(wd, os.Stderr)
			if err != nil {
				return err
			}
			return repo.Fetch(os.Stdin, os.Stdout)
		},
	}
}

func NewPullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "fetch chunks for split files in the working tree and combine",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, _ := os.Getwd()
			repo, err := bits.NewRepository(wd, os.Stderr)
			if err != nil {
				return err
			}
			return repo.Pull("HEAD", os.Stdout)
		},
	}
}

func NewPushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "push locally stored chunks to the remote store",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, _ := os.Getwd()
			repo, err := bits.NewRepository(wd, os.Stderr)
			if err != nil {
				return err
			}
			store, err := repo.LocalStore()
			if err != nil {
				return err
			}
			defer store.Close()
			return repo.Push(store, os.Stdin, "origin")
		},
	}
}

func NewCombineCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "combine",
		Short: "combine chunks back into the original file",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, _ := os.Getwd()
			repo, err := bits.NewRepository(wd, os.Stderr)
			if err != nil {
				return err
			}
			return repo.Combine(os.Stdin, os.Stdout)
		},
	}
}