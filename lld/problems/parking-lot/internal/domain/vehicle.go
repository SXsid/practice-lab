package domain

type VehicalType string

const (
	Car        VehicalType = "car"
	Truck      VehicalType = "truck"
	MotorCycle VehicalType = "motorcycle"
)

type Vehical struct {
	Length      int
	Width       int
	Type        VehicalType
	Id          string
	NumberPlate string
	TicketId    string
}
type Vehical interface {
	Type() VehicalType
	Dimensions() (int, int)
	NumberPlate() string
	Ticket() string
}
