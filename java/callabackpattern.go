package main

// Good method fuction callback

import (
	"fmt"
	"time"
)

// funiton clalback=>stor funtio in the main struc instef of sving the refer of c hild struc  who has the fuciton
// it's like single chain chefk take an d build now we have to reigst which custoemr just thign without creating hte cusotmer
// instace
type (
	callback func(dish string)
	chefMain struct {
		CallbackFun callback
	}
)

func NewChefMain() *chefMain {
	return &chefMain{
		CallbackFun: func(dish string) {
			fmt.Printf("here you go sir you %s,anything else", dish)
		},
	}
}

func (c *chefMain) PlaceOrder(dish string) {
	fmt.Println("ok")
	fmt.Println("preapring")
	fmt.Println("done")
	// now the custoe which order mean which invoked the New chek
	c.CallbackFun(dish)
}

// java styule or secoed way
type servingStaff interface {
	placeOrder(dish string)
	onReady(dish string)
}

type Staff struct{}

type chef struct {
	staffToserver servingStaff
}

func NewChef(staff servingStaff) *chef {
	return &chef{
		staff,
	}
}

func (c *chef) prepareDish(dish string) {
	fmt.Println("preparin dish ")
	time.Sleep(2 * time.Second)
	fmt.Printf("%s is prepared\n", dish)
	c.staffToserver.onReady(dish)
}

func (s *Staff) placeOrder(dish string) {
	chef := NewChef(s)
	chef.prepareDish(dish)
}

func (s *Staff) onReady(dish string) {
	fmt.Println("here is your dish sir anything else")
}

func run() {
	NewChefMain().PlaceOrder("chili panner")
	Staff := Staff{}
	Staff.placeOrder("dal makhni")
}
