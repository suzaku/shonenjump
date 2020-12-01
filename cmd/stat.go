package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/suzaku/shonenjump/jump"
)

func init() {
	rootCmd.AddCommand(statCmd)
}

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "Show information about recorded paths",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		dataPath := jump.EnsureDataPath()
		entries := jump.LoadEntries(dataPath)
		for _, e := range entries {
			fmt.Println(e)
		}
	},
}
