package data

type CommandCard struct {
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Requirements []string `json:"requirements"`
	Faction      string   `json:"faction"`
	Keywords     []string `json:"keywords"`
	Pips         int      `json:"pips"`
	Commander    string   `json:"commander"`
	Orders       string   `json:"orders"`
	Weapon       *Weapon  `json:"weapon,omitempty"`
	Text         string   `json:"text,omitempty"`
}
