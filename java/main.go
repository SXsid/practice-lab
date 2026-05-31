package main

import "fmt"

type NotifcationProvider interface {
	send(msg string)
}
type NotifcationServie struct {
	provider NotifcationProvider
}

func NewNotifcatService(provider NotifcationProvider) *NotifcationServie {
	return &NotifcationServie{
		provider: provider,
	}
}

type sms struct{}

func (s *sms) send(msg string) {
	fmt.Println("sms sedning")
}

type email struct{}

func (e *email) send(msg string) {
	fmt.Println("email sedning")
}

func main() {
	NewNotifcatService(&email{}).provider.send("hi")
	NewNotifcatService(&sms{}).provider.send("hi")
}
