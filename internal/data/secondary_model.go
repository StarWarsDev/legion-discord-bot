package data

import (
	"fmt"
	"strings"
)

type Dice struct {
	Black int `json:"black,omitempty"`
	Red   int `json:"red,omitempty"`
	White int `json:"white,omitempty"`
}

func (dice *Dice) String() string {
	var strSlice []string

	if dice.White > 0 {
		strSlice = append(strSlice, fmt.Sprintf("white: %d", dice.White))
	}

	if dice.Black > 0 {
		strSlice = append(strSlice, fmt.Sprintf("black: %d", dice.Black))
	}

	if dice.Red > 0 {
		strSlice = append(strSlice, fmt.Sprintf("red: %d", dice.Red))
	}

	return strings.Join(strSlice, ", ")
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

func (weapon *Weapon) String() string {
	var weaponInfo []string
	if len(weapon.Name) > 0 {
		weaponInfo = append(weaponInfo, "  "+weapon.Name)
	}

	var to interface{}
	to = weapon.Range.To
	// this looks like it is an infinite range weapon
	if weapon.Range.To < weapon.Range.From {
		to = "∞"
	}
	weaponInfo = append(weaponInfo, fmt.Sprintf("Range: %d - %v", weapon.Range.From, to))
	weaponInfo = append(weaponInfo, fmt.Sprintf("Dice: %s", weapon.Dice.String()))

	if len(weapon.Keywords) > 0 {
		keywords := strings.Join(weapon.Keywords, ", ")
		weaponInfo = append(weaponInfo, "Keywords: "+keywords)
	}

	return strings.Join(weaponInfo, "\n")
}
