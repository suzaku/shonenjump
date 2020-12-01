package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/suzaku/shonenjump/jump"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add this path",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := jump.AddPath(args[0]); err != nil {
			log.Fatal(err)
		}
	},
}
