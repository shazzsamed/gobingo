package game

import (
	"encoding/json"
	"fmt"
	"gobingo/helpers"
	"gobingo/models"
	"log"
	"net/url"
	"strconv"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

var WssURL string 
type ServerMessage struct {
	Event string      `json:"event"`
	Data  json.RawMessage `json:"data"`
}

var conn *websocket.Conn
var myTurn bool
var roomCode string

func PlayWithFriend() {
	fmt.Println("Enter a name to start the game:")
	var name string
	fmt.Scan(&name)

	fmt.Println("Choose an option: \n1. Create a Room\n2. Join a Room")
	var choice int
	fmt.Scan(&choice)

	var roomCode string

	switch choice {
	case 1:
		roomCode = createRoom(name)
	case 2:
		joinRoom(name)
	default:
		fmt.Println("Invalid choice.")
		return
	}

	board := setupBoard(roomCode)
	sendBoardReady(name, board)
	waitForGameStart(name)
}

func createRoom(name string) string {
	var err error
	u := url.URL{
		Scheme: "wss",
		Host: WssURL,
		Path: "/create-room", 
		RawQuery: "name=" + name,
	}

	conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}

	var response map[string]string
	if err := conn.ReadJSON(&response); err != nil{
		log.Fatal("Error reading response:", err)
	}
	roomCode = response["roomCode"]
	return roomCode
}

func joinRoom(name string) {
	var err error

	fmt.Print("Enter Room Code: ")
	var joinCode string
	fmt.Scan(&joinCode)

	u := url.URL{
		Scheme: "wss",
		Host: WssURL,
		Path: "/join-room",
		RawQuery: "roomCode=" + joinCode + "&name=" + name,
	}
	conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}

    var response map[string]string
    if err := conn.ReadJSON(&response); err != nil {
        log.Fatal("Error reading response from server:", err)
    }

    if errorMessage, exists := response["error"]; exists {
        fmt.Println("Error:", errorMessage)
        return
    }
}

func setupBoard(roomCode string) [][]models.Cell {
	board := helpers.Create2DArray(5, 5, -1)
	displayError := ""
	usedNumbers := make(map[int]bool)
	markedNumbers := make(map[int]bool)
	var input string
	OUTERLOOP:
    for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			helpers.ClearTerminal()
            printBoard(board, markedNumbers)
			fmt.Println("Share this roomCode with your friend: ", roomCode)
            color.Red(displayError)
            fmt.Printf("(enter 'random' or 'r' to generate random bingo card)\n")
            fmt.Printf("Enter value for cell [%d, %d]: ", i+1, j+1)
            fmt.Scan(&input)
            if input == "r" || input == "random" {
                board = helpers.CreateRandomBoard()
                break OUTERLOOP
            }
            num, err := strconv.Atoi(input)
            if err != nil || num < 1 || num > 25{
                j--
                displayError = "Invalid input. Please enter a number between 1 and 25."
                continue
            }

            if usedNumbers[num] {
                j--
                displayError = "Number already used. Please enter different number."
                continue
            }
            board[i][j].Number = num
            board[i][j].Marked = false
            usedNumbers[num] = true
            displayError = ""
		}
	}
	helpers.ClearTerminal()
	printBoard(board, markedNumbers)
	fmt.Println("Waiting for the other player to setup the board...")
	fmt.Println("Room Code: ", roomCode)
	return board
}

func sendBoardReady(playerName string, board [][]models.Cell) {
	message := map[string]interface{}{
		"action": "boardReady",
		"playerName" : playerName,
		"board": board,
	}
	if err := conn.WriteJSON(message); err != nil {
        log.Fatal("Error sending board ready:", err)
	}
}

func waitForGameStart(name string) {
	for {
		var msg ServerMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Fatal("Error reading server message:", err)
		}

		switch msg.Event {
		case "startGame":
			var data map[string]interface{}
			if err := json.Unmarshal(msg.Data, &data); err != nil {
				log.Println("Failed to unmarshal startGame data:", err)
				continue
			}

			if firstPlayer, ok := data["nextTurn"].(string); ok {
				myTurn = (firstPlayer == name)
			}

			playGame(name)
			return
		}
	}
}

func playGame(name string) {
	for {
        if myTurn {
            fmt.Println("Your Turn! Enter a number to mark:")
            var num float64
            fmt.Scan(&num)

            message := map[string]interface{}{
                "action": "markNumber",
                "number": num,
            }
            conn.WriteJSON(message)
            data, err := processServerMessage()
			if err != nil {
				log.Println("Error:", err)
			} else {
				if (data.Bingo!="none"){
					fmt.Println("THE WINNER IS: ", data.Bingo)
					fmt.Println("Enter any key to exit")
					var exit string
					fmt.Scan(&exit)
					break
				}
				printBothBoard(data.YourBoard, data.OppBoard, bingo, 0)
				fmt.Println("Marked number: ",data.Number,", Next Turn", data.NextTurn)
				if data.NextTurn == name {
					myTurn = true
				} else {
					myTurn = false
				}
			}
        } else {
            data, err := processServerMessage()
			if err != nil {
				log.Println("Error:", err)
			} else {
				if (data.Bingo!="none"){
					fmt.Println("THE WINNER IS: ", data.Bingo)
					fmt.Println("Enter any key to exit")
					var exit string
					fmt.Scan(&exit)
					break
				}
				printBothBoard(data.YourBoard, data.OppBoard, bingo, 0)
				fmt.Println("Marked number: ",data.Number,", Next Turn", data.NextTurn)
				if data.NextTurn == name {
					myTurn = true
				} else {
					myTurn = false
				}
			}
			continue
        }
	}
}

func processServerMessage() (models.MarkNumberData, error) {
	var msg ServerMessage
	var markData models.MarkNumberData

	if err := conn.ReadJSON(&msg); err != nil {
		return markData, fmt.Errorf("error reading server message: %v", err)
	}

	switch msg.Event {
	case "markNumber":
		if err := json.Unmarshal(msg.Data, &markData); err != nil {
			return markData, fmt.Errorf("invalid markNumber data: %v", err)
		}
		return markData, nil

	case "Error":
		var errMsg string
		if err := json.Unmarshal(msg.Data, &errMsg); err != nil {
			return markData, fmt.Errorf("unknown server error")
		}
		return markData, fmt.Errorf(errMsg)

	default:
		return markData, fmt.Errorf("unhandled event type: %v", msg.Event)
	}
}
