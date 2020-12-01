package cmd

import (
	"github.com/spf13/cobra"

	"github.com/suzaku/shonenjump/jump"
)

func init() {
	rootCmd.AddCommand(purgeCmd)
}

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Remove non-existent paths from database",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		dataPath := jump.EnsureDataPath()
		entries := jump.LoadEntries(dataPath)
		entries, changed := jump.ClearNotExistDirs(entries)
		if changed {
			entries.Save(dataPath)
		}
	},
}
