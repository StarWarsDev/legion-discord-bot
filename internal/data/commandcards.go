package data

type CommandCard struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Icon         string   `json:"icon"`
	Image        string   `json:"image"`
	Requirements []string `json:"requirements"`
	Faction      string   `json:"faction"`
	Keywords     []string `json:"keywords"`
	Pips         int      `json:"pips"`
	Commander    string   `json:"commander"`
	Orders       string   `json:"orders"`
	Weapon       *Weapon  `json:"weapon,omitempty"`
}
