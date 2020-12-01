package cmd

import (
	"fmt"

	"github.com/suzaku/shonenjump/jump"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(guessCmd)
}

var guessCmd = &cobra.Command{
	Use:   "guess",
	Short: "Output a best guess of path for the argument",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dataPath := jump.EnsureDataPath()
		if len(args) == 1 {
			needle, index, path := jump.ParseCompleteOption(args[0])
			if path != "" {
				fmt.Println(path)
				return
			} else if index != 0 {
				path = jump.GetNCandidate([]string{needle}, index, ".")
				fmt.Println(path)
				return
			} else {
				args = []string{needle}
			}
		}
		entries := jump.LoadEntries(dataPath)
		fmt.Println(jump.BestGuess(entries, args))
	},
}
