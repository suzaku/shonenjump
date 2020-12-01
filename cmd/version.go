package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/suzaku/shonenjump/jump"
)

func init() {
	rootCmd.AddCommand(verCmd)
}

var verCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version of shonenjump",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(jump.VERSION)
	},
}
