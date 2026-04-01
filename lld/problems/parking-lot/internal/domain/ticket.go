package domain

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	id       string
	spot     *Spot
	vehical  *Vehical
	starTime time.Time
	endTime  time.Time
	ownerId  string
	amount   float64
}

// type Ticket interface {
// 	VerifyTicket(ownerId string) bool
// Amount() float64
// 	VehicalDetail() Vehical
// 	SpotDetail() Spot
// TimeitTook()
// }

func NewTicket(spot *Spot, vehcial *Vehical, ownerId string) *Ticket {
	return &Ticket{
		id:       uuid.NewString(),
		spot:     spot,
		vehical:  vehcial,
		starTime: time.Now(),
		ownerId:  ownerId,
	}
}

func (t *Ticket) VerifyTicket(ownerId string) bool {
	return t.ownerId == ownerId
}

func (t *Ticket) Amount() float64 {
	t.endTime = time.Now()
	return  t.spot.Price()*
}
