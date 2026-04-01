package service

import "github.com/SXsid/lld-problems/parking-lot/internal/domain"

type BookingService interface {
	Park(Vehical domain.Vehical) (string, error)
	PayableAmount(ticketID, ownerId string) (float64, error)
	Unpark()
}
