package cmd

import (
	"fmt"
	"gobingo/assets"
	"gobingo/game"
	"math/rand"
	"time"

	"github.com/inancgumus/screen"
	"github.com/spf13/cobra"
)

var (
	playWithComputer bool
	playWithFriend   bool
)

func init() {
	playCmd.Flags().BoolVarP(&playWithComputer, "computer", "c", false, "Play against the computer")
	playCmd.Flags().BoolVarP(&playWithFriend, "friend", "f", false, "Play with a friend")

	rootCmd.AddCommand(playCmd)
}

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Start a Bingo game",
	Long:  "Start a Bingo game. Choose to play against a computer or a friend using flags.",
	Run: func(cmd *cobra.Command, args []string) {
		screen.Clear()
		screen.MoveTopLeft()
		assets.TitleCard()
		if playWithComputer && playWithFriend {
			fmt.Println("Error: Please select only one option: --computer or --friend.")
			return
		}

		if playWithComputer {
			playBingo(true)
		} else if playWithFriend {
			playBingo(false)
		} else {
			fmt.Println("Choose your game mode:\n1. Play with Gobi (Computer)\n2. Play with a friend")
			var choice int
			fmt.Print("Enter your choice (1 or 2): ")
			fmt.Scan(&choice)
			if choice==1 {
				game.PlayWithComputer()
			} else if choice==2 {
				game.PlayWithFriend()
			}
		}
	},
}

func playBingo(vsComputer bool) {
	board := generateBoard()
	fmt.Println("Your Bingo Board:")
	printBoard(board)

	if vsComputer {
		fmt.Println("Playing against the computer...")
		// Add game logic against the computer
	} else {
		fmt.Println("Playing with a friend...")
		// Add game logic for friends
	}
}

func generateBoard() [5][5]int {
	rand.Seed(time.Now().UnixNano())
	nums := rand.Perm(25) // Randomly shuffle numbers from 0 to 24
	var board [5][5]int
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			board[i][j] = nums[i*5+j] + 1 // Map 0-24 to 1-25
		}
	}
	return board
}

func printBoard(board [5][5]int) {
	for _, row := range board {
		for _, num := range row {
			fmt.Printf("%2d ", num)
		}
		fmt.Println()
	}
}
