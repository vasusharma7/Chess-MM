/*
Contains server code for go chess server
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var gameCache map[string]Board = make(map[string]Board)

// MoveRequestBody received to move a piece
type MoveRequestBody struct {
	FromRow int `json:"FromRow"`
	ToRow   int `json:"ToRow"`
	FromCol int `json:"FromCol"`
	ToCol   int `json:"ToCol"`
}

// MoveResponseBody sent as result
type MoveResponseBody struct {
	Board [][]string `json:"Board"`
	Check bool       `json:"Check"`
	Mate  bool       `json:"Mate"`
}

func play(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(r.URL.Path)
	if r.URL.Path != "/engine/api/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	var (
		oldPos, newPos, fromPos, toPos Position
		score                          float64
		allPos                         []Position
	)
	id := r.URL.Query().Get("id")
	fmt.Println(id)
	board, ok := gameCache[id]
	if !ok {
		fmt.Println("Initialising a new game...")
		board = Board{}
		board.initialise()
		gameCache[id] = board
	}
	board.print()

	switch r.Method {
	case "GET":

		// jData, err := json.Marshal(board.getAsSlice())
		// if err != nil {
		// 	fmt.Fprintf(w, "json.Marshal() err: %v", err)
		// 	return
		// }
		var res MoveResponseBody
		res.Board = board.getAsSlice()
		res.Check = false
		res.Mate = false
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)

	case "POST":
		var body MoveRequestBody
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Could not decode request data", body)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.Unmarshal(buf, &body)
		fmt.Println(body)

		fromPos, toPos = Position{body.FromRow, body.FromCol}, Position{body.ToRow, body.ToCol}

		allPos, err = board[fromPos.row][fromPos.col].getAllMoves(board, fromPos)

		if err != nil {
			fmt.Println("No element found at ", fromPos)
			http.Error(w, "No element found", http.StatusNonAuthoritativeInfo)
			return
		}

		if !contains(allPos, toPos) {
			fmt.Println("Not a valid move for element at position:", fromPos)
			http.Error(w, "Not a valid move", http.StatusNotAcceptable)
			return
		}

		if First(board.movePiece(fromPos, toPos)).check(User) == 1 {
			fmt.Println("You cannot move here, your king will be in check position")
			http.Error(w, "Cannot move here, King will be in CHECK state", http.StatusForbidden)
			return
		}
		board.makeMove(fromPos, toPos)

		//check for stalemate by generating all moves
		fmt.Println("Hmm....nice move....you have forced me to hit my nerves...")
		oldPos, newPos, score = miniMax(0, Tree{board: board}, Self, MIN, MAX)
		fmt.Println(oldPos, newPos, score)

		if score == MIN {
			fmt.Println("Looks like no moves left for me !")
			http.Error(w, "No moves left for me :( ", http.StatusNoContent)
			return
		}

		board.makeMove(oldPos, newPos)
		gameCache[id] = board
		board.print()
		var res MoveResponseBody
		res.Board = board.getAsSlice()
		res.Check = false
		res.Mate = false

		if board.check(User) == 1 {
			fmt.Println("CHECK !")
			res.Check = true
		}

		if len(board.generateNodes(User)) == 0 {
			fmt.Println("OR Looks like CHECK AND MATE !")
			res.Mate = true
		}

		// jData, err := json.Marshal(res)
		// if err != nil {
		// 	fmt.Fprintf(w, "json.Marshal() err: %v", err)
		// 	return
		// }
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	http.HandleFunc("/", play)

	fmt.Printf("Starting Chess Server...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func (board Board) getAsSlice() (boardList [][]string) {
	for i := 0; i < len(board); i++ {
		rowList := []string{}
		for j := 0; j < len(board[i]); j++ {
			rowList = append(rowList, board[i][j].String())
		}
		boardList = append(boardList, rowList)

	}
	return
}

func formBoardUsingSlice(boardList [][]string) (board Board) {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			board[i][j] = &Empty{}
		}
	}
	for i := 0; i < len(boardList); i++ {
		for j := 0; j < len(boardList[i]); j++ {
			if boardList[i][j] != (&Empty{}).String() {
				board[i][j] = getPieceFromString(boardList[i][j])
			}
		}

	}
	return
}

func getPieceFromString(pStr string) Piece {
	switch pStr {
	case "K":
		return &King{User}
	case "K'":
		return &King{Self}
	case "Q":
		return &Queen{User}
	case "Q'":
		return &Queen{Self}
	case "B":
		return &Bishop{User}
	case "B'":
		return &Bishop{Self}
	case "N":
		return &Knight{User}
	case "N'":
		return &Knight{Self}
	case "P":
		return &Pawn{User}
	case "P'":
		return &Pawn{Self}
	case "R":
		return &Rook{User}
	case "R'":
		return &Rook{Self}

	default:
		return &Empty{}
	}
}

func contains(positions []Position, loc Position) bool {
	for _, v := range positions {
		if v == loc {
			return true
		}
	}
	return false
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
