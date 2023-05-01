/*
This is executable package which will run the chess engine.
*/
package main

import (
	"fmt"
)

// overall goal - make attacking from defensive
func main() {
	board := Board{}
	board.initialise()
	// for i := 0; i < 8; i++ {
	// 	for j := 0; j < 8; j++ {
	// 		board[i][j] = &Empty{}
	// 	}
	// }
	// board[0][1] = &Rook{Self}
	// board[0][2] = &Pawn{User}
	// board[1][1] = &Pawn{User}
	// board[1][5] = &King{Self}
	// board[2][2] = &Pawn{User}

	// board[2][3] = &Pawn{Self}
	// board[2][7] = &Knight{Self}

	// board[3][0] = &Pawn{Self}
	// board[3][1] = &Knight{Self}

	// board[4][4] = &Pawn{User}
	// board[4][7] = &Pawn{User}

	// board[5][3] = &King{User}

	board.print()
	var (
		from, to                       string
		oldPos, newPos, fromPos, toPos Position
		score                          float64
		allPos                         []Position
		err                            error
	)
	for {
		fmt.Printf("move from: ")
		fmt.Scanf("%s \n", &from)
		if !validateInput(from) {
			fmt.Println("Invalid Input")
			continue
		}
		fmt.Printf("move to: ")
		fmt.Scanf("%s \n", &to)
		if !validateInput(to) {
			fmt.Println("Invalid Input")
			continue
		}
		fmt.Println(from, "->", to)

		fromPos = getPositionFromInput(from)
		toPos = getPositionFromInput(to)

		allPos, err = board[fromPos.row][fromPos.col].getAllMoves(board, fromPos)
		if err != nil {
			fmt.Println("No element found at ", from)
			continue
		}
		if !contains(allPos, toPos) {
			fmt.Println("Not a valid move for element at position:", from)
			continue
		}
		if First(board.movePiece(fromPos, toPos)).check(User) == 1 {
			fmt.Println("You cannot move here, your king will be in check position")
			continue
		}
		board.makeMove(fromPos, toPos)
		board.print()
		//check for stalemate by generating all moves
		fmt.Println("Hmm....nice move....you have forced me to hit my nerves...")
		oldPos, newPos, score = miniMax(0, Tree{board: board}, Self, MIN, MAX)
		if score == MAX {
			fmt.Println("Looks like no moves left for me !")
			break
		}
		board.makeMove(oldPos, newPos)
		board.print()
		if board.check(User) == 1 {
			println("CHECK !")
		}
		if len(board.generateNodes(User)) == 0 {
			fmt.Println("OR Looks like CHECK AND MATE !")
			break
		}
		fmt.Println(oldPos, newPos, score)
	}
}

func getPositionFromInput(input string) Position {
	return Position{7 - (int(input[1]) - 49), int(input[0]) - 97}
}

func validateInput(input string) bool {
	if len(input) != 2 {
		return false
	}
	if int(input[0]) < 97 || int(input[0]) > 104 {
		return false
	}
	if int(input[1]) < 49 || int(input[1]) > 56 {
		return false
	}
	return true
}

func contains(positions []Position, loc Position) bool {
	for _, v := range positions {
		if v == loc {
			return true
		}
	}
	return false
}
