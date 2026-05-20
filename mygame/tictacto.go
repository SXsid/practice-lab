package main

import "fmt"

type TicTacToe struct {
	board [3][3]string
}

func NewTickTacToe() *TicTacToe {
	board := [3][3]string{
		{"_", "_", "_"},
		{"_", "_", "_"},
		{"_", "_", "_"},
	}
	return &TicTacToe{
		board: board,
	}
}

func (t *TicTacToe) PrintGrid() {
	for i := range len(t.board) {
		for j := range len(t.board[0]) {
			fmt.Print(t.board[i][j])
		}
		fmt.Println()
	}
}
