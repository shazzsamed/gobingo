package helpers

import (
	"fmt"
	"gobingo/models"
	"time"

	"github.com/inancgumus/screen"
	"golang.org/x/exp/rand"
)

func CreateRandomBoard() [][]models.Cell {
	board := Create2DArray(5,5, -1)
	sequence := generateShuffledSequence(25)
	for row := 0; row < 5; row++ {
		for col := 0; col < 5; col++ {
			board[row][col].Number = sequence[(row*5)+col]
			board[row][col].Marked = false
		}
	}
	return board
}

func generateShuffledSequence(n int) []int {
	sequence := make([]int, n)
	for i := 0; i < n; i++ {
		sequence[i] = i + 1
	}
	rand.Seed(uint64(time.Now().UnixNano()))
	rand.Shuffle(len(sequence), func(i, j int) {
		sequence[i], sequence[j] = sequence[j], sequence[i]
	})
	return sequence
}

func Create2DArray(rows, cols int, defaultValue int) [][]models.Cell {
	array := make([][]models.Cell, rows)
	for i := range array {
		array[i] = make([]models.Cell, cols)
		for j := range array[i] {
			array[i][j].Number = defaultValue
			array[i][j].Marked = false
		}
	}
	return array
}

func Create2DIntArray(rows, cols int, defaultValue int) [][]int {
	board := make([][]int, rows)
	for i := range board {
		board[i] = make([]int, cols)
		for j := range board[i] {
			board[i][j] = defaultValue
		}
	}
	return board
}

func ClearTerminal() {
	screen.MoveTopLeft()
	screen.Clear()
	// fmt.Print("\033[H\033[2J")
	fmt.Print("\033c") // Reset the terminal and clear everything, including scrollback
}