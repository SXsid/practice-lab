package service

import (
	"fmt"

	"github.com/SXsid/lld-problems/parking-lot/internal/domain"
	"github.com/google/uuid"
)

type BookingService interface {
	GetVacantSpots(vehical domain.Vehical) []domain.Spot
	PayableAmount(ticketID, ownerId string) (float64, error)
	ProcessPayment(ticketId string) error
	TicketById(ticketId string) *domain.Ticket
	ListTicket(limit, offset int) []domain.Ticket
	Park(vehical *domain.Vehical, spot *domain.Spot) (string, error)
}

type TicketRepo interface {
	Add(ticket domain.Ticket) error
	GetByID(id string) *domain.Ticket
	List(limt, offset int) []domain.Ticket
	Update(ticket domain.Ticket, id string) error
}
type SpotRepo interface {
	Add(spot domain.Spot) error
	GetByID(id string) *domain.Spot
	List(limt, offset int) []domain.Spot
	Update(spot domain.Spot, id string) error
}
type PaymentRepo interface {
	pay(amount float64) error
}
type Booking struct {
	// dummy repo
	spotRepo    SpotRepo
	ticketRepo  TicketRepo
	paymentRepo PaymentRepo
}

func NewBooking(spotRepo SpotRepo, ticketRepo TicketRepo, paymentRepo PaymentRepo) *Booking {
	return &Booking{
		spotRepo:    spotRepo,
		ticketRepo:  ticketRepo,
		paymentRepo: paymentRepo,
	}
}

func (b *Booking) TicketById(ticketId string) (*domain.Ticket, error) {
	return b.ticketRepo.GetByID(ticketId), nil
}

func (b *Booking) ListTicket(limit, offset int) []domain.Ticket {
	return b.ticketRepo.List(limit, offset)
}

func (b *Booking) PayableAmount(ticketID, ownerId string) (float64, error) {
	t, err := b.TicketById(ticketID)
	if err != nil {
		return 0, err
	}
	if !t.VerifyTicket(ownerId) {
		return 0, fmt.Errorf("Invalid ownerId")
	}
	if t.Status() == domain.Active {
		t.GeneratBill()
	}
	return t.Amount(), nil
}

// after razorpay or the pahment interface webook call it
func (b *Booking) ProcessPayment(ticketId string) error {
	t, err := b.TicketById(ticketId)
	if err != nil {
		return err
	}
	t.Paid()
	if err := b.ticketRepo.Update(*t, ticketId); err != nil {
		return err
	}
	t.Close()
	if err := b.ticketRepo.Update(*t, ticketId); err != nil {
		return err
	}

	spot := t.Spot()
	if err := b.spotRepo.Update(*spot, spot.Id()); err != nil {
		return err
	}

	return nil
}

func (b *Booking) GetVacantSpots(VehticalType domain.VehicalType) ([]domain.Spot, error) {
	// idk how to implemt db here like shoul di store flor or the just the spot and make floor strc out of it
	Spots := []domain.Spot{}
	if len(Spots) < 1 {
		return nil, fmt.Errorf("no spot avaiblbl")
	}
	return Spots, nil
}

func (b *Booking) Park(vehical *domain.Vehical, spot *domain.Spot) (string, error) {
	if !spot.IsCompatible(vehical.Type()) || spot.Status() {
		return "", fmt.Errorf("spot is not compatible with vehcial or alredy allocated ")
	}

	ticketId := uuid.NewString()
	ticket := domain.NewTicket(ticketId, spot, vehical)

	spot.ChangeStatus(true)
	b.ticketRepo.Add(*ticket)

	return ticketId, nil
}
