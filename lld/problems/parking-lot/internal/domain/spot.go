package domain

type SpotType string

const (
	CarSpot        SpotType = "car"
	TruckSpot      SpotType = "truck"
	MotorCycleSpot SpotType = "motorcycle"
)

type Spot struct {
	id          string
	spotType    SpotType
	isAllocated bool
	lenght      int
	width       int
	price       float64
}

// type Spot interface {
// 	Status() bool
// 	ChangeStatus(isAvaiblabe bool)
// 	Dimension() (int, int)
// 	IsCompatible(vehicalType VehicalType) bool
// }

func NewSpot(spotType SpotType, price float64, l, b int) *Spot {
	return &Spot{
		spotType:    spotType,
		lenght:      l,
		width:       b,
		isAllocated: false,
		price:       price,
	}
}

func (s *Spot) Price() float64 {
	return s.price
}

func (s *Spot) Status() bool {
	return s.isAllocated
}

func (s *Spot) ChangeStatus(isAvailable bool) {
	s.isAllocated = isAvailable
}

func (s *Spot) Dimension() (int, int) {
	return s.lenght, s.width
}

func (s *Spot) IsCompatible(vehicleType VehicalType) bool {
	switch s.spotType {
	case CarSpot:
		return vehicleType == Car
	case TruckSpot:
		return vehicleType == Truck
	case MotorCycleSpot:
		return vehicleType == MotorCycle
	default:
		return false
	}
}

func (s *Spot) Id() string {
	return s.id
}
