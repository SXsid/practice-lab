package domain

type Owner struct {
	idCard   string
	name     string
	vehicals []Vehical
}

// type Owner interface {
// 	IdCard() string
// 	Vehicals() []Vehical
//  AddVehicals(vehcial Vehical)
// }

func NewOwner(name, idCard string) *Owner {
	return &Owner{
		idCard:   idCard,
		name:     name,
		vehicals: make([]Vehical, 0),
	}
}

func (o *Owner) IdCard() string {
	return o.idCard
}

func (o *Owner) AddVehical(vehical Vehical) {
	o.vehicals = append(o.vehicals, vehical)
}

func (o *Owner) Vehicals() []Vehical {
	return o.vehicals
}
