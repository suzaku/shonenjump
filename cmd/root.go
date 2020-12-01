package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "shonenjump",
		Short: "A faster way to change directory and improve command line productivity.",
	}
)

func Execute() error {
	return rootCmd.Execute()
}
