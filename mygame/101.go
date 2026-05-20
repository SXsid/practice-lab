package main

import (
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/manifoldco/promptui"
)

type Level string

const (
	EASY   Level = "easy"
	MEDIUM Level = "medium"
	HARD   Level = "hard"
)

var ErrInvalidType = errors.New("invalid game difficulty type")

func playGame(lives int) {
	score := 0
	fmt.Println(score)
	for {
		num := Pick()
		if num == rand.IntN(3)+1 {
			score++
		} else {
			lives--
		}
		fmt.Printf("(score:%d,lifeRemaining:%d)\n", score, lives)
		if score == 5 {
			fmt.Println("You win!")
			break
		}
		// after score as the life got cosmed to print/ calcu this loop
		if lives == 0 {
			fmt.Println("Game Over")
			break
		}

	}
}

func Pick() int {
	var num int
	fmt.Println("Enter anumber 1 out of 3")
	fmt.Scan(&num)
	return num
}

func printOption() int {
	items := []string{string(EASY), string(MEDIUM), string(HARD)}
	prompt := promptui.Select{
		Label: "Select a difficulty",
		Items: items,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	switch Level(result) {
	case EASY:
		return 5
	case MEDIUM:
		return 3
	case HARD:
		return 1
	default:
		panic(ErrInvalidType.Error())
	}
}

func start101() {
	fmt.Println("Game Started")
	lives := printOption()
	playGame(lives)
}
