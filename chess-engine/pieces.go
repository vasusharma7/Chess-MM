package main

import (
	"errors"
)

// Piece that encapsulates methods on board.
type Piece interface {
	getAllMoves(Board, Position) ([]Position, error)

	String() string

	getPlayer() Color
}

// Empty board state
type Empty struct{}

// Pawn piece
type Pawn struct {
	color Color
}

// King piece
type King struct {
	color Color
}

// Queen Piece
type Queen struct {
	color Color
}

// Rook piece
type Rook struct {
	color Color
}

// Bishop piece
type Bishop struct {
	color Color
}

// Knight piece
type Knight struct {
	color Color
}

func (Empty) String() string {
	return "-"
}

func (p Knight) String() string {
	if p.color == Self {
		return "N'"
	}
	return "N"
}

func (p King) String() string {
	if p.color == Self {
		return "K'"
	}
	return "K"
}

func (p Queen) String() string {
	if p.color == Self {
		return "Q'"
	}
	return "Q"
}

func (p Rook) String() string {
	if p.color == Self {
		return "R'"
	}
	return "R"
}

func (p Bishop) String() string {
	if p.color == Self {
		return "B'"
	}
	return "B"
}

func (p Pawn) String() string {
	if p.color == Self {
		return "P'"
	}
	return "P"
}

func (p Empty) getPlayer() Color {
	return Undefined
}

func (p Knight) getPlayer() Color {
	return p.color
}

func (p King) getPlayer() Color {
	return p.color
}

func (p Queen) getPlayer() Color {
	return p.color
}

func (p Rook) getPlayer() Color {
	return p.color
}

func (p Bishop) getPlayer() Color {
	return p.color
}

func (p Pawn) getPlayer() Color {
	return p.color
}

func (p Empty) getAllMoves(board Board, loc Position) ([]Position, error) {
	return []Position{}, nil
}

func (p Pawn) getAllMoves(board Board, loc Position) ([]Position, error) {
	pos := []Position{}
	if board[loc.row][loc.col].getPlayer() == Undefined {
		return nil, errors.New("Piece not present at the location")
	}

	if board[loc.row][loc.col].getPlayer() == Self {

		if loc.row+1 < 8 && board[loc.row+1][loc.col].getPlayer() == Undefined {
			pos = append(pos, Position{loc.row + 1, loc.col})
		}

		if loc.row+1 < 8 && loc.col+1 < 8 && board[loc.row+1][loc.col+1].getPlayer() == User {
			pos = append(pos, Position{loc.row + 1, loc.col + 1})
		}

		if loc.row+1 < 8 && loc.col-1 >= 0 && board[loc.row+1][loc.col-1].getPlayer() == User {
			pos = append(pos, Position{loc.row + 1, loc.col - 1})
		}
	} else {
		if loc.row-1 >= 0 && board[loc.row-1][loc.col].getPlayer() == Undefined {
			pos = append(pos, Position{loc.row - 1, loc.col})
		}

		if loc.row-1 >= 0 && loc.col+1 < 8 && board[loc.row-1][loc.col+1].getPlayer() == Self {
			pos = append(pos, Position{loc.row - 1, loc.col + 1})
		}

		if loc.row-1 >= 0 && loc.col-1 >= 0 && board[loc.row-1][loc.col-1].getPlayer() == Self {
			pos = append(pos, Position{loc.row - 1, loc.col - 1})
		}
	}
	return pos, nil
}

func (p Rook) getAllMoves(board Board, loc Position) ([]Position, error) {
	pos := []Position{}
	pos = append(pos, board.getFwMoves(loc)...)
	pos = append(pos, board.getBkMoves(loc)...)
	pos = append(pos, board.getLtMoves(loc)...)
	pos = append(pos, board.getRtMoves(loc)...)

	return pos, nil
}

func (p Bishop) getAllMoves(board Board, loc Position) ([]Position, error) {
	pos := []Position{}
	pos = append(pos, board.getBLDiagonal(loc)...)
	pos = append(pos, board.getFRDiagonal(loc)...)
	pos = append(pos, board.getBRDiagonal(loc)...)
	pos = append(pos, board.getFLDiagonal(loc)...)

	return pos, nil
}

func (p Queen) getAllMoves(board Board, loc Position) ([]Position, error) {
	pos := []Position{}
	pos = append(pos, board.getBLDiagonal(loc)...)
	pos = append(pos, board.getFRDiagonal(loc)...)
	pos = append(pos, board.getBRDiagonal(loc)...)
	pos = append(pos, board.getFLDiagonal(loc)...)
	pos = append(pos, board.getFwMoves(loc)...)
	pos = append(pos, board.getBkMoves(loc)...)
	pos = append(pos, board.getLtMoves(loc)...)
	pos = append(pos, board.getRtMoves(loc)...)

	return pos, nil
}

func (p King) getAllMoves(board Board, loc Position) ([]Position, error) {
	pos := []Position{}
	var against Color
	if board[loc.row][loc.col].getPlayer() == Undefined {
		return pos, errors.New("Piece not present at the location")
	}

	if board[loc.row][loc.col].getPlayer() == White {
		against = Black
	} else {
		against = White
	}

	if loc.row+1 < 8 && board.isValidMove(against, Position{loc.row + 1, loc.col}) {
		pos = append(pos, Position{loc.row + 1, loc.col})
	}

	if loc.row-1 >= 0 && board.isValidMove(against, Position{loc.row - 1, loc.col}) {
		pos = append(pos, Position{loc.row - 1, loc.col})
	}

	if loc.col+1 < 8 && board.isValidMove(against, Position{loc.row, loc.col + 1}) {
		pos = append(pos, Position{loc.row, loc.col + 1})
	}

	if loc.col-1 >= 0 && board.isValidMove(against, Position{loc.row, loc.col - 1}) {
		pos = append(pos, Position{loc.row, loc.col - 1})
	}

	if loc.row+1 < 8 && loc.col+1 < 8 && board.isValidMove(against, Position{loc.row + 1, loc.col + 1}) {
		pos = append(pos, Position{loc.row + 1, loc.col + 1})
	}

	if loc.row-1 >= 0 && loc.col-1 >= 0 && board.isValidMove(against, Position{loc.row - 1, loc.col - 1}) {
		pos = append(pos, Position{loc.row - 1, loc.col - 1})
	}

	if loc.row+1 < 8 && loc.col-1 >= 0 && board.isValidMove(against, Position{loc.row + 1, loc.col - 1}) {
		pos = append(pos, Position{loc.row + 1, loc.col - 1})
	}

	if loc.row-1 >= 0 && loc.col+1 < 8 && board.isValidMove(against, Position{loc.row - 1, loc.col + 1}) {
		pos = append(pos, Position{loc.row - 1, loc.col + 1})
	}
	return pos, nil
}

func (p Knight) getAllMoves(board Board, loc Position) ([]Position, error) {
	pos := []Position{}
	var against Color
	if board[loc.row][loc.col].getPlayer() == White {
		against = Black
	} else {
		against = White
	}
	if loc.row+2 < 8 && loc.col+1 < 8 && board.isValidMove(against, Position{loc.row + 2, loc.col + 1}) {
		pos = append(pos, Position{loc.row + 2, loc.col + 1})
	}

	if loc.row+2 < 8 && loc.col-1 >= 0 && board.isValidMove(against, Position{loc.row + 2, loc.col - 1}) {
		pos = append(pos, Position{loc.row + 2, loc.col - 1})
	}

	if loc.row-2 >= 0 && loc.col+1 < 8 && board.isValidMove(against, Position{loc.row - 2, loc.col + 1}) {
		pos = append(pos, Position{loc.row - 2, loc.col + 1})
	}

	if loc.row-2 >= 0 && loc.col-1 >= 0 && board.isValidMove(against, Position{loc.row - 2, loc.col - 1}) {
		pos = append(pos, Position{loc.row - 2, loc.col - 1})
	}

	if loc.row+1 < 8 && loc.col+2 < 8 && board.isValidMove(against, Position{loc.row + 1, loc.col + 2}) {
		pos = append(pos, Position{loc.row + 1, loc.col + 2})
	}

	if loc.row+1 < 8 && loc.col-2 >= 0 && board.isValidMove(against, Position{loc.row + 1, loc.col - 2}) {
		pos = append(pos, Position{loc.row + 1, loc.col - 2})
	}

	if loc.row-1 >= 0 && loc.col+2 < 8 && board.isValidMove(against, Position{loc.row - 1, loc.col + 2}) {
		pos = append(pos, Position{loc.row - 1, loc.col + 2})
	}

	if loc.row-1 >= 0 && loc.col-2 >= 0 && board.isValidMove(against, Position{loc.row - 1, loc.col - 2}) {
		pos = append(pos, Position{loc.row - 1, loc.col - 2})
	}

	return pos, nil
}

func (board Board) isValidMove(against Color, pos Position) bool {
	return board[pos.row][pos.col].getPlayer() == Undefined || board[pos.row][pos.col].getPlayer() == against
}
