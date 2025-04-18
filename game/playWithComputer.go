package game

import (
	"fmt"
	"gobingo/helpers"
	"gobingo/models"
	"strconv"
	"time"

	"github.com/fatih/color"
	"golang.org/x/exp/rand"
)

var Header = "╔════╦════╦════╦════╦════╗"
var Footer = "╚════╩════╩════╩════╩════╝"
var Break = "╠════╬════╬════╬════╬════╣"
var Line = "║%-4v║%-4v║%-4v║%-4v║%-4v║"
var bingo = [5]string{"B","I","N","G","O"}
var green = color.New(color.FgGreen).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()

// PlayWithComputer initializes the board and handles interactive input.
// allow user to genrste board (do not repeat numbers, d for delete, c for confirm after all done)
func PlayWithComputer() {
    board := helpers.Create2DArray(5,5, -1)
    usedNumbers := make(map[int]bool)
    markedNumbers := make(map[int]bool)
    var input string
    displayError := ""
    OuterLoop:
    for row := 0; row < 5; row++ {
        for col := 0; col < 5; col++ {
            helpers.ClearTerminal()
            printBoard(board,markedNumbers)
            color.Red(displayError)
            fmt.Printf("(enter 'random' or 'r' to generate random bingo card)\n")
            fmt.Printf("Enter value for cell [%d, %d]: ", row+1, col+1)
            fmt.Scan(&input)
            if input == "r" || input == "random" {
                board = helpers.CreateRandomBoard()
                break OuterLoop
            }
            num, err := strconv.Atoi(input)
            if err != nil || num < 1 || num > 25{
                col-- 
                displayError = "Invalid input. Please enter a number between 1 and 25."
                continue
            }

            if usedNumbers[num] {
                col--
                displayError = "Number already used. Please enter different number."
                continue
            }
            board[row][col].Number = num
            board[row][col].Marked = false
			usedNumbers[num] = true
            displayError = ""
        }
    }
    helpers.ClearTerminal()
    oppBoard := helpers.CreateRandomBoard()
    startGame(board, oppBoard)
}

// printBoard dynamically displays the Bingo table.
func printBoard(board [][]models.Cell, marked map[int]bool) {
	fmt.Println(Header)
	for i, row := range board {
		var formattedRow [5]interface{}
		for j, num := range row {
			// fmt.Printf("num: %v\n", num.Number)
            if num.Number == -1 {
                formattedRow[j] = "   "
            } else if num.Marked {
				formattedRow[j] = green(num.Number) // Print in green if marked
			} else {
				formattedRow[j] = num.Number
			}
		}
		fmt.Printf(Line+"\n", formattedRow[:]...)
		if i < len(board)-1 {
			fmt.Println(Break)
		}
	}
	fmt.Println(Footer)
}

func printBothBoard(yourBoard, oppBoard [][]models.Cell, bingo [5]string, count int) {
	helpers.ClearTerminal()
	fmt.Println(Header + "    " + Header) // Print headers for both boards

	for i := 0; i < 5; i++ {
		// Format player board
		var formattedYourRow [5]string
        for j, cell := range yourBoard[i] {
            if cell.Number == -1 {
                formattedYourRow[j] = "   "
            } else if cell.Marked {
                formattedYourRow[j] = green(fmt.Sprintf("%-4d", cell.Number))
            } else {
                formattedYourRow[j] = fmt.Sprintf("%-4d", cell.Number)
            }
        }

        var formattedOppRow [5]string
        for j, cell := range oppBoard[i] {
            if (cell.Number == -1 && !cell.Marked) {
                formattedOppRow[j] = "   "
            } else if cell.Marked {
                formattedOppRow[j] = green("X   ")
            } else {
                formattedOppRow[j] = "    " 
            }
        }

        // Now print correctly with all values as strings
        fmt.Printf(Line, formattedYourRow[0], formattedYourRow[1], formattedYourRow[2], formattedYourRow[3], formattedYourRow[4])
        if i < count {
			fmt.Printf(cyan("%-4s"), bingo[i])
		} else {
			fmt.Printf("%-4s", bingo[i])
		}
        fmt.Printf(Line+"\n", formattedOppRow[0], formattedOppRow[1], formattedOppRow[2], formattedOppRow[3], formattedOppRow[4])


		// Print separator if not the last row
		if i < 4 {
			fmt.Println(Break + "    " + Break)
		}
	}
	fmt.Println(Footer + "    " + Footer) // Print footers for both boards
}

func startGame(playerBoard, oppBoard [][]models.Cell) {
	markedNumbers := make(map[int]bool)
	var input string
    var countBingo int

	for {
		// Print both boards
		printBothBoard(playerBoard, oppBoard, bingo, countBingo)

		// Get user input
		fmt.Print("Enter a number to mark: ")
		fmt.Scan(&input)
		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > 25 {
			fmt.Println("Invalid input. Please enter a number between 1 and 25.")
			continue
		}

		// Check if the number is already marked
		if markedNumbers[num] {
			fmt.Println("Number already marked. Please enter a different number.")
			continue
		}

		// Mark the player's board
		markedNumbers[num] = true
		markBoard(playerBoard, num)
		markBoard(oppBoard, num)

        isBingo, countBingo := checkBingo(playerBoard)
		// Check for Bingo
		if isBingo {
			fmt.Println("Congratulations! You won! BINGO!")
			return
		}

		printBothBoard(playerBoard, oppBoard, bingo, countBingo)

        // Opponent's turn
		oppNum := selectOpponentMove(playerBoard)
		fmt.Printf("Opponent chose: %d\n", oppNum)

		// Mark the boards for the opponent's move
		markBoard(playerBoard, oppNum)
		markBoard(oppBoard, oppNum)

        isBingoOpp, _ := checkBingo(playerBoard)
		// Check for Bingo
		if isBingoOpp {
			fmt.Println("Opponent won! BINGO!")
			return
		}
	}
}

func markBoard(board [][]models.Cell, num int) {
    for i := 0; i < 5; i++ {
        for j := 0; j < 5; j++ {
            if board[i][j].Number == num {
                board[i][j].Marked = true
                board[i][j].Number = num
                return
            }
        }
    }
}

func checkBingo(board [][]models.Cell) (bool, int) {
	count := 0
	for i := 0; i < 5; i++ {
		if isRowMarked(board, i) {
			count++
		}
	}
	for j := 0; j < 5; j++ {
		if isColumnMarked(board, j) {
			count++
		}
	}
	if isMainDiagonalMarked(board) {
		count++
	}
	return count >= 5, count
}

func isRowMarked(board [][]models.Cell, row int) bool {
	for j := 0; j < 5; j++ {
		if !board[row][j].Marked {
			return false
		}
	}
	return true
}

func isColumnMarked(board [][]models.Cell, col int) bool {
	for i := 0; i < 5; i++ {
		if !board[i][col].Marked {
			return false
		}
	}
	return true
}

func isMainDiagonalMarked(board [][]models.Cell) bool {
	for i := 0; i < 5; i++ {
		if !board[i][i].Marked {
			return false
		}
	}
	return true
}

func selectOpponentMove(oppBoard [][]models.Cell) int {
    delay := time.Duration(rand.Intn(2)+3) * time.Second // 3-7 seconds delay
    time.Sleep(delay)

    for {
        num := rand.Intn(25) + 1 // Random number between 1-25

        if !isNumberMarked(oppBoard, num) {
            return num
        }
    }
}

// Helper function to check if a number is already marked
func isNumberMarked(board [][]models.Cell, num int) bool {
    for i := 0; i < 5; i++ {
        for j := 0; j < 5; j++ {
            if board[i][j].Number == num && board[i][j].Marked {
                return true
            }
        }
    }
    return false
}
