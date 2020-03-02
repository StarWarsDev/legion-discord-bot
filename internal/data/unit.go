package data

// Unit is a standard Legion Unit
type Unit struct {
	Name         string        `json:"name"`
	Image        string        `json:"image"`
	Requirements []string      `json:"requirements"`
	Faction      string        `json:"faction"`
	Keywords     []string      `json:"keywords"`
	Unique       bool          `json:"unique"`
	Cost         int           `json:"cost"`
	Rank         string        `json:"rank"`
	Slots        []string      `json:"slots"`
	CommandCards []CommandCard `json:"commandCards"`
	Wounds       int           `json:"wounds"`
	Courage      *int          `json:"courage,omitempty"`
	Resilience   *int          `json:"resilience,omitempty"`
	Defense      string        `json:"defense"`
	Entourage    []string      `json:"entourage"`
	Surge        *Surge        `json:"surge,omitempty"`
	Weapons      []Weapon      `json:"weapons"`
	Type         string        `json:"cardSubType" graphql:"cardSubType"`
}

func (unit *Unit) FullName() string {
	name := unit.Name
	if unit.Unique {
		name = "• " + name
	}
	return name
}
