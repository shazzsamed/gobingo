package models

type Cell struct {
	Number int
	Marked bool
}

type MarkNumberData struct {
	YourBoard [][]Cell `json:"yourBoard"`
	OppBoard  [][]Cell `json:"oppBoard"`
	Number    int      `json:"number"`
	Player    string   `json:"player"`
	NextTurn  string   `json:"nextTurn"`
	Bingo     string   `json:"bingo"`
}