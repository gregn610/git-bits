package command

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
	"github.com/nerdalize/git-bits/bits"
)

func NewInstallCmd() *cobra.Command {
	var bucket, remote string
	
	cmd := &cobra.Command{
		Use:   "install",
		Short: "configures filters, create pre-push hook and pull chunks",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %v", err)
			}

			repo, err := bits.NewRepository(wd, os.Stderr)
			if err != nil {
				return fmt.Errorf("failed to setup repository: %v", err)
			}

			conf := bits.DefaultConf()
			
			if bucket == "" {
				conf.AWSS3BucketName, err = askInput("In which AWS S3 bucket would you like to store chunks? ")
				if err != nil {
					return fmt.Errorf("failed to get bucket input: %v", err)
				}
			} else {
				conf.AWSS3BucketName = bucket
			}

			conf.AWSAccessKeyID, err = askInput("What is your AWS Access Key ID? ")
			if err != nil {
				return fmt.Errorf("failed to get access key input: %v", err)
			}

			conf.AWSSecretAccessKey, err = askSecret("What is your AWS Secret Key? (input will be hidden) ")
			if err != nil {
				return fmt.Errorf("failed to get secret key input: %v", err)
			}

			err = repo.Install(os.Stdout, conf)
			if err != nil {
				return fmt.Errorf("failed to install: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&bucket, "bucket", "b", "", "name of the s3 bucket used as a chunk remote")
	cmd.Flags().StringVarP(&remote, "remote", "r", "origin", "git remote that will be configured for chunk storage")

	return cmd
}

func askInput(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func askSecret(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println()
	return string(bytePassword), nil
}
