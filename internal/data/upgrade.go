package data

type Upgrade struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Text               string   `json:"text"`
	Type               string   `json:"type" graphql:"cardSubType"`
	Image              string   `json:"image"`
	Requirements       []string `json:"requirements"`
	Keywords           []string `json:"keywords"`
	Cost               int      `json:"cost"`
	Exhaust            bool     `json:"exhaust"`
	Weapon             *Weapon  `json:"weapon"`
	UnitTypeExclusions []string `json:"unitTypeExclusions"`
}
