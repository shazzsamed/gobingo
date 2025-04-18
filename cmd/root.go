package cmd

import (
	"fmt"
	"gobingo/assets"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bingo",
	Short: "Bingo CLI Game",
	Long:  "A fancy 5x5 Bingo CLI game where you can play with a computer or a friend.",
	Run: func(cmd *cobra.Command, args []string) {
		assets.TitleCard()
		fmt.Println("Use the 'play' command to start a game. Use --help for more details.")
	},
}

// Execute initializes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
