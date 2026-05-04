// Copyright (c) 2026 Keith Chu
package cmd

import (
	"github.com/cqroot/minbox/pkg/version"
	"github.com/spf13/cobra"
)

// runRootCmd is the default root command that executes all operations.
func runRootCmd(cmd *cobra.Command, args []string) {
	// Do Stuff Here
}

// newRootCmd creates and returns the root cobra command.
func newRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "minbox",
		Short: "minbox - A CLI tool",
		Long:  `minbox is a CLI tool built with Cobra.`,
		Run:   runRootCmd,
	}
	rootCmd.AddCommand(newBaseCmd())
	rootCmd.Version = version.Get().String()
	return &rootCmd
}

func Execute() {
	cobra.CheckErr(newRootCmd().Execute())
}
