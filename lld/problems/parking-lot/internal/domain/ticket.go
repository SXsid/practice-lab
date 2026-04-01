package domain

import (
	"time"

	"github.com/sethvargo/go-retry"
)

type TicketStatus string

const (
	Active TicketStatus = "active"
	Unpaid TicketStatus = "unpaid"
	Paid   TicketStatus = "paid"
	Close  TicketStatus = "close"
)

type Ticket struct {
	id       string
	spot     *Spot
	vehical  *Vehical
	starTime time.Time
	endTime  time.Time
	amount   float64
	status   TicketStatus
}

func NewTicket(id string, spot *Spot, vehcial *Vehical) *Ticket {
	return &Ticket{
		id:       id,
		spot:     spot,
		vehical:  vehcial,
		starTime: time.Now(),
		status:   Active,
	}
}

func (t *Ticket) VerifyTicket(ownerId string) bool {
	return t.vehical.OwnerId() == ownerId
}

func (t *Ticket) Status() TicketStatus {
	return t.status
}

func (t *Ticket) Amount() float64 {
	return t.amount
}

func (t *Ticket) GeneratBill() {
	end := time.Now()
	t.endTime = end
	diff := end.Sub(t.starTime).Hours()
	amount := t.spot.Price() * diff
	t.amount = amount
	t.status = Unpaid
}

func (t *Ticket) Paid() {
	t.status = Paid
}

func (t *Ticket) Close() {
	t.status = Close
	t.spot.ChangeStatus(true)
}

func (t *Ticket) Spot() *Spot {
	return t.spot
}
