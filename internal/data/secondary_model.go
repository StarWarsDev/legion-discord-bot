package data

type Dice struct {
	Black int `json:"black,omitempty"`
	Red   int `json:"red,omitempty"`
	White int `json:"white,omitempty"`
}

type Range struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type Surge struct {
	Attack  string `json:"attack,omitempty"`
	Defense string `json:"defense,omitempty"`
}

type Weapon struct {
	Name     string   `json:"name,omitempty"`
	Range    Range    `json:"range"`
	Dice     Dice     `json:"dice"`
	Surge    *Surge   `json:"surge,omitempty"`
	Keywords []string `json:"keywords"`
}
