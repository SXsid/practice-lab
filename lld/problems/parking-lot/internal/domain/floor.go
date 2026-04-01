package domain

type FloorStruct struct {
	Number int
	Spots  []Spot
}
type Floor interface {
	GetData() ([]Floor, error)
	GetAvaibleSlotFloorVise(floorNumber int) ([]domain.Spot, error)
	GetAvaibleSlotForAVechicalType(vechincalType domain.VehicalType) ([]domain.Spot, error)
}
