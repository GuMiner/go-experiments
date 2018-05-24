package citizen

type Citizen struct {
	Age int // In days
}

type Population struct {
	citizens      map[int64]*Citizen
	nextCitizenId int64
}

func newCitizen() *Citizen {
	return &Citizen{}
}

func NewPopulation() *Population {
	return &Population{
		citizens:      make(map[int64]*Citizen),
		nextCitizenId: 0}
}
