package cmd

import (
	"github.com/spf13/cobra"

	"github.com/suzaku/shonenjump/jump"
)

func init() {
	rootCmd.AddCommand(completeCmd)
}

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Used for tab completion",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var arg string
		if len(args) > 0 {
			arg = args[0]
		}
		jump.ShowAutoCompleteOptions(arg)
	},
}
