/*
Contains chess board related methods and types.
*/
package main

import (
	"fmt"
)

// Color of piece
type Color int

const (
	// White piece
	White Color = iota
	// Black piece
	Black
	//Undefined or No piece
	Undefined
)

// Board that represents chess game
type Board [8][8]Piece

// Position on board of form a1, b4 etc...
type Position struct {
	row int `json:"row,omitempty"`
	col int `json:"col,omitempty"`
}

func getColor(position Position) Color {
	if (position.row+position.col)%2 == 0 {
		return White
	}
	return Black
}

func (board *Board) initialise() {
	board[0][0] = &Rook{Self}
	board[0][1] = &Knight{Self}
	board[0][2] = &Bishop{Self}
	board[0][3] = &Queen{Self}
	board[0][4] = &King{Self}
	board[0][5] = &Bishop{Self}
	board[0][6] = &Knight{Self}
	board[0][7] = &Rook{Self}

	for i := 0; i < 8; i++ {
		board[1][i] = &Pawn{Self}
		board[6][i] = &Pawn{User}
	}

	board[7][0] = &Rook{User}
	board[7][1] = &Knight{User}
	board[7][2] = &Bishop{User}
	board[7][3] = &Queen{User}
	board[7][4] = &King{User}
	board[7][5] = &Bishop{User}
	board[7][6] = &Knight{User}
	board[7][7] = &Rook{User}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board[i][j] == nil {
				board[i][j] = &Empty{}
			}
		}
	}
}

func (board Board) print() {
	//fmt.Print("\033[H\033[2J")
	fmt.Println()
	for i := 0; i < len(board); i++ {
		fmt.Printf("%d |\t", 8-i)
		for j := 0; j < len(board[i]); j++ {
			fmt.Print(board[i][j])
			if board[i][j].getPlayer() == Self {
				fmt.Print("\t")
			} else {
				fmt.Print(" \t")
			}
		}
		fmt.Println()
		fmt.Print("  |\n")
	}
	fmt.Printf("   ")
	for i := 0; i < len(board)*4; i++ {
		fmt.Print("__")
	}
	fmt.Println()
	fmt.Println()
	fmt.Printf(" \t")
	for i := 0; i < len(board); i++ {
		fmt.Printf("%c \t", 'a'+i)
	}
	fmt.Println()
}

// MovePiece moves piece on board from old to new position and returns a new board
func (board Board) movePiece(oldPosition Position, newPosition Position) (Board, error) {
	board[oldPosition.row][oldPosition.col], board[newPosition.row][newPosition.col] = &Empty{}, board[oldPosition.row][oldPosition.col]
	return board, nil
}

// makeMove moves piece on board from old to new position
func (board *Board) makeMove(oldPosition Position, newPosition Position) error {
	board[oldPosition.row][oldPosition.col], board[newPosition.row][newPosition.col] = &Empty{}, board[oldPosition.row][oldPosition.col]
	return nil
}

// GetAllHorizontalMoves ...
func GetAllHorizontalMoves(position Position) (hPosn []Position) {
	for i := 0; i < 8; i++ {
		hPosn = append(hPosn, Position{position.row, i})
	}
	return
}

func (board Board) getFwMoves(loc Position) (pos []Position) {
	self := board[loc.row][loc.col].getPlayer()
	for i := loc.row + 1; i < 8; i++ {
		if board[i][loc.col].getPlayer() == self {
			break
		}
		pos = append(pos, Position{i, loc.col})

		if board[i][loc.col].getPlayer() != Undefined {
			break
		}
	}
	return pos
}

func (board Board) getBkMoves(loc Position) (pos []Position) {
	self := board[loc.row][loc.col].getPlayer()
	for i := loc.row - 1; i >= 0; i-- {
		if board[i][loc.col].getPlayer() == self {
			break
		}
		pos = append(pos, Position{i, loc.col})

		if board[i][loc.col].getPlayer() != Undefined {
			break
		}
	}
	return pos
}

func (board Board) getRtMoves(loc Position) (pos []Position) {
	self := board[loc.row][loc.col].getPlayer()
	for i := loc.col + 1; i < 8; i++ {
		if board[loc.row][i].getPlayer() == self {
			break
		}
		pos = append(pos, Position{loc.row, i})

		if board[loc.row][i].getPlayer() != Undefined {
			break
		}
	}
	return pos
}

func (board Board) getLtMoves(loc Position) (pos []Position) {
	self := board[loc.row][loc.col].getPlayer()
	for i := loc.col - 1; i >= 0; i-- {
		if board[loc.row][i].getPlayer() == self {
			break
		}
		pos = append(pos, Position{loc.row, i})

		if board[loc.row][i].getPlayer() != Undefined {
			break
		}
	}
	return pos
}

// GetAllVerticalMoves ...
func getAllVerticalMoves(position Position) (vPos []Position) {
	for i := 0; i < 8; i++ {
		vPos = append(vPos, Position{i, position.col})
	}
	return
}

// GetFRDiagonal ...
func (board Board) getFRDiagonal(position Position) (dPos []Position) {
	i := position.row
	j := position.col
	self := board[i][j].getPlayer()

	for {
		i++
		j++
		if i > 7 || j > 7 || board[i][j].getPlayer() == self {
			break
		}

		dPos = append(dPos, Position{i, j})
		if board[i][j].getPlayer() != Undefined {
			break
		}
	}
	return dPos
}

// GetFLDiagonal ...
func (board Board) getFLDiagonal(position Position) (dPos []Position) {
	i := position.row
	j := position.col
	self := board[i][j].getPlayer()

	for {
		i++
		j--
		if i > 7 || j < 0 || (board[i][j].getPlayer() == self) {
			break
		}

		dPos = append(dPos, Position{i, j})
		if board[i][j].getPlayer() != Undefined {
			break
		}
	}
	return dPos
}

// getBLDiagonal ...
func (board Board) getBLDiagonal(position Position) (dPos []Position) {
	i := position.row
	j := position.col
	self := board[i][j].getPlayer()
	for {
		i--
		j--
		if i < 0 || j < 0 || (board[i][j].getPlayer() == self) {
			break
		}

		dPos = append(dPos, Position{i, j})
		if board[i][j].getPlayer() != Undefined {
			break
		}
	}
	return
}

// getBRDiagonal ...
func (board Board) getBRDiagonal(position Position) (dPos []Position) {
	i := position.row
	j := position.col
	self := board[i][j].getPlayer()

	for {
		i--
		j++
		if i < 0 || j > 7 || board[i][j].getPlayer() == self {
			break
		}

		dPos = append(dPos, Position{i, j})
		if board[i][j].getPlayer() != Undefined {
			break
		}
	}
	return dPos
}
