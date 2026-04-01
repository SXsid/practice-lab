package domain

type VehicalType string

const (
	Car        VehicalType = "car"
	Truck      VehicalType = "truck"
	MotorCycle VehicalType = "motorcycle"
)

type Vehical struct {
	length      int
	width       int
	vehicalType VehicalType
	id          string
	numberPlate string
}

// type Vehical interface {
// 	Type() VehicalType
// 	Dimensions() (int, int)
// 	NumberPlate() string
// 	Ticket() string
// }

func NewVehical(owner *Owner, l, b int, numberplate string, vehicalType VehicalType) *Vehical {
	return &Vehical{
		length:      l,
		width:       b,
		id:          owner.IdCard(),
		numberPlate: numberplate,
		vehicalType: vehicalType,
	}
}

func (v *Vehical) Type() VehicalType {
	return v.vehicalType
}

func (v *Vehical) OwnerId() string {
	return v.id
}
