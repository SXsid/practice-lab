package main

import "fmt"

type Institute struct{}

func (i Institute) Teach() {
	fmt.Println("Teaches")
}

type College struct {
	Institute
}

func (c College) Sports() {
	fmt.Println("College Sports")
}

func (c College) Reserach() {
	fmt.Println("Research college")
}

// compositotn over inherticae  abstraic achive by iinteface
type Uni struct {
	College
}

func (u Uni) Reserach() {
	fmt.Println("uni Research college")
}

func Run() {
	uni := Uni{}
	uni.Teach()
	uni.Reserach()
	uni.Sports()
}
