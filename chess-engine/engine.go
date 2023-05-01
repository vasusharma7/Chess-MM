/*
Contains chess engine related methods and types.
*/
package main

import (
	"errors"
	"math"
	"math/rand"
	"reflect"
)

// TODO: Check and protect from CHECK to King.
const (
	//Self Color
	Self Color = Black
	// User Piece Color
	User Color = White
	// MaxDepth of MiniMax Tree
	MaxDepth = 4
	// MAX number
	MAX = float64(1000)
	//MIN number
	MIN = float64(-1000)
)

// Tree structure for minimax algorithm
type Tree struct {
	board  Board
	oldPos Position
	newPos Position
	nodes  []Tree
	score  int
}

/** Zorbist hashing here*/

// var cache [MaxDepth]map[int]Tree

// var zorbistTable [8][8][12]int
// func hash(board Board) int {

// }

func miniMax(depth int, tree Tree, player Color,
	alpha float64, beta float64) (oldPos Position, newPos Position, score float64) {
	if depth == MaxDepth {
		return tree.oldPos, tree.newPos, tree.board.evaluate()
	}

	if player == Self { //maximizer
		best := MIN
		index := -1
		tree.nodes = append(tree.nodes, tree.board.generateNodes(Self)...)
		for i := 0; i < len(tree.nodes); i++ {
			_, _, val := miniMax(depth+1, tree.nodes[i], User, alpha, beta)
			//experimental - for risk taking
			if First(tree.board.movePiece(tree.nodes[i].oldPos, tree.nodes[i].newPos)).check(User) == 1 && val > alpha  {
				best = val
				alpha = math.Max(best, alpha)
				index = i
				break
			}
			alpha = math.Max(alpha, val)
			if best < val {
				best = val
				index = i
			}
			// best = math.Max(best, val)
			if beta <= alpha {
				break
			}
		}
		if index == -1 {
			return tree.oldPos, tree.newPos, MAX
		}
		return tree.nodes[index].oldPos, tree.nodes[index].newPos, best
	}
	best := MAX
	index := -1
	tree.nodes = append(tree.nodes, tree.board.generateNodes(User)...)
	for i := 0; i < len(tree.nodes); i++ {
		_, _, val := miniMax(depth+1, tree.nodes[i], Self, alpha, beta)
		//experimental - for risk taking
		if First(tree.board.movePiece(tree.nodes[i].oldPos, tree.nodes[i].newPos)).check(Self) == 1 && val < beta {
			best = val
			beta = math.Min(best, beta)
			index = i
			break
		}

		beta = math.Min(beta, val)
		if best > val {
			best = val
			index = i
		}
		if beta <= alpha {
			break
		}
	}
	if index == -1 {
		return tree.oldPos, tree.newPos, MIN
	}
	return tree.nodes[index].oldPos, tree.nodes[index].newPos, best
}

func (board Board) generateNodes(color Color) []Tree {
	nodes := []Tree{}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board[i][j].getPlayer() == color {
				moves, _ := board[i][j].getAllMoves(board, Position{i, j})
				for _, move := range moves {
					newBoard := First(board.movePiece(Position{i, j}, move))
					if newBoard.check(color) == 1 {
						continue
					}
					nodes = append(nodes, Tree{board: newBoard, oldPos: Position{i, j}, newPos: move})
				}
			}
		}
	}
	return shuffle(nodes)
}

// enhance this algo
// use blocked pawns + isolated pawns + doubled pawns
func (board Board) evaluate() float64 {
	kingWt := 200.0 * float64(board.getPieceDiff(reflect.TypeOf(&King{})))
	queenWt := 9.0 * float64(board.getPieceDiff(reflect.TypeOf(&Queen{})))
	rookWt := 5.0 * float64(board.getPieceDiff(reflect.TypeOf(&Rook{})))
	bishopWt := 3.0 * float64(board.getPieceDiff(reflect.TypeOf(&Bishop{})))
	knightWt := 3.0 * float64(board.getPieceDiff(reflect.TypeOf(&Knight{})))
	pawnWt := 1.0 * float64(board.getPieceDiff(reflect.TypeOf(&Pawn{})))
	checkWt := 100.0 * float64((board.check(User) - board.check(Self)))
	mobility := 0.1 * float64(board.check(Self)-board.check(User))
	return float64(kingWt + queenWt + rookWt + bishopWt + knightWt + pawnWt + mobility + checkWt)
}

func (board Board) check(player Color) int {
	var opponent Color
	if player == User {
		opponent = Self
	} else {
		opponent = User
	}
	pos, err := board.findPiece(King{player})
	if err != nil {
		return 1
	}

	blDiagonalPos := getLastElement(board.getBLDiagonal(pos))
	if board.checkDialogalDanger(blDiagonalPos, opponent) {
		return 1
	}

	brDiagonalPos := getLastElement(board.getBRDiagonal(pos))
	if board.checkDialogalDanger(brDiagonalPos, opponent) {
		return 1
	}

	flDiagonalPos := getLastElement(board.getFLDiagonal(pos))
	if board.checkDialogalDanger(flDiagonalPos, opponent) {
		return 1
	}

	frDiagonalPos := getLastElement(board.getFRDiagonal(pos))
	if board.checkDialogalDanger(frDiagonalPos, opponent) {
		return 1
	}

	ltPos := getLastElement(board.getLtMoves(pos))
	if board.checkAxialDanger(ltPos, opponent) {
		return 1
	}

	fwdPos := getLastElement(board.getFwMoves(pos))
	if board.checkAxialDanger(fwdPos, opponent) {
		return 1
	}

	bkPos := getLastElement(board.getBkMoves(pos))
	if board.checkAxialDanger(bkPos, opponent) {
		return 1
	}

	rtPos := getLastElement(board.getRtMoves(pos))
	if board.checkAxialDanger(rtPos, opponent) {
		return 1
	}

	knightMoves, _ := Knight{}.getAllMoves(board, pos)
	for _, move := range knightMoves {
		if board[move.row][move.col].String() == (&Knight{opponent}).String() {
			return 1
		}
	}

	kingMoves, _ := King{}.getAllMoves(board, pos)
	for _, move := range kingMoves {
		if board[move.row][move.col].String() == (&King{opponent}).String() {
			return 1
		}
	}

	//Pawns
	if player == Self {
		if pos.row+1 < 8 && pos.col-1 >= 0 && board[pos.row+1][pos.col-1].String() == (&Pawn{User}).String() {
			return 1
		}
		if pos.row+1 < 8 && pos.col+1 < 8 && board[pos.row+1][pos.col+1].String() == (&Pawn{User}).String() {
			return 1
		}
	} else {
		if pos.row-1 >= 0 && pos.col-1 >= 0 && board[pos.row-1][pos.col-1].String() == (&Pawn{Self}).String() {
			return 1
		}
		if pos.row-1 >= 0 && pos.col+1 < 8 && board[pos.row-1][pos.col+1].String() == (&Pawn{Self}).String() {
			return 1
		}
	}

	return 0
}

func (board Board) checkDialogalDanger(position Position, opponent Color) bool {
	if position.row == -1 || position.col == -1 {
		return false
	}
	if board[position.row][position.col].getPlayer() == opponent {
		piece := reflect.TypeOf(board[position.row][position.col])
		if piece == reflect.TypeOf(&Queen{}) || piece == reflect.TypeOf(&Bishop{}) {
			return true
		}
	}
	return false
}

func (board Board) checkAxialDanger(position Position, opponent Color) bool {
	if position.row == -1 || position.col == -1 {
		return false
	}
	if board[position.row][position.col].getPlayer() == opponent {
		piece := reflect.TypeOf(board[position.row][position.col])
		if piece == reflect.TypeOf(&Queen{}) || piece == reflect.TypeOf(&Rook{}) {
			return true
		}
	}
	return false
}

func getLastElement(positions []Position) Position {
	if len(positions) == 0 {
		return Position{-1, -1}
	}
	return positions[len(positions)-1]
}

func (board Board) findPiece(piece Piece) (Position, error) {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board[i][j].String() == piece.String() {
				return Position{i, j}, nil
			}
		}
	}
	return Position{-1, -1}, errors.New("Piece not found")
}

func (board Board) getPieceDiff(id reflect.Type) int {
	return board.countTypeOfPiece(id, Self) - board.countTypeOfPiece(id, User)
}

func (board Board) countTypeOfPiece(id reflect.Type, player Color) int {
	count := 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board[i][j].getPlayer() == player && reflect.TypeOf(board[i][j]) == id {
				count++
			}
		}
	}
	return count
}

// First returns first element from return values
func First(b Board, _ error) Board {
	return b
}

// shuffle shuffles the elements of an array in place
func shuffle(array []Tree) []Tree {
	for i := range array { //run the loop till the range of array
		j := rand.Intn(i + 1)                   //choose any random number
		array[i], array[j] = array[j], array[i] //swap the random element with current element
	}
	return array
}
