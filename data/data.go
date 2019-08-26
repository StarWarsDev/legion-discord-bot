package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// LegionData is a container for all Legion data
type LegionData struct {
	Units        Units
	Upgrades     Upgrades
	CommandCards []CommandCard
}

// Units is a container of all the different unit types
type Units struct {
	Commander     []Unit
	Operative     []Unit
	Corps         []Unit
	SpecialForces []Unit `json:"Special Forces"`
	Support       []Unit
	Heavy         []Unit
}

// Flattened returns a flat slice of units regardless of type
func (u *Units) Flattened() (units []Unit) {
	units = append(units, u.Commander...)
	units = append(units, u.Operative...)
	units = append(units, u.Corps...)
	units = append(units, u.SpecialForces...)
	units = append(units, u.Support...)
	units = append(units, u.Heavy...)

	return
}

// Unit is a Legion unit card
type Unit struct {
	LDF          string `json:"ldf"`
	Name         string
	Subtitle     string
	Type         string
	Points       int
	Rank         string
	Minis        int
	Wounds       int
	Courage      int
	Resilience   int
	Defense      string
	Surge        Surge
	Speed        int
	Slots        []string
	Keywords     []Keyword
	Weapons      []Weapon
	CommandCards []string
	Unique       bool
}

// Surge denotes which kinds of surges a unit has
type Surge struct {
	Attack  string
	Defense string
}

// Upgrades is a container of all the different upgrade types
type Upgrades struct {
	Armament    []Upgrade
	Command     []Upgrade
	Comms       []Upgrade
	Elite       []Upgrade
	Force       []Upgrade
	Gear        []Upgrade
	Generator   []Upgrade
	Grenades    []Upgrade
	Hardpoint   []Upgrade
	HeavyWeapon []Upgrade `json:"Heavy Weapon"`
	Personnel   []Upgrade
	Pilot       []Upgrade
}

// Flattened returns a flat slice of all upgrades regardless of type
func (u *Upgrades) Flattened() (upgrades []Upgrade) {
	upgrades = append(upgrades, u.Armament...)
	upgrades = append(upgrades, u.Command...)
	upgrades = append(upgrades, u.Comms...)
	upgrades = append(upgrades, u.Elite...)
	upgrades = append(upgrades, u.Force...)
	upgrades = append(upgrades, u.Gear...)
	upgrades = append(upgrades, u.Generator...)
	upgrades = append(upgrades, u.Grenades...)
	upgrades = append(upgrades, u.Hardpoint...)
	upgrades = append(upgrades, u.HeavyWeapon...)
	upgrades = append(upgrades, u.Personnel...)
	upgrades = append(upgrades, u.Pilot...)

	return
}

// Upgrade is a Legion upgrade card
type Upgrade struct {
	LDF          string `json:"ldf"`
	Name         string
	Description  string
	Points       int
	Restrictions UpgradeRestrictions
	Slot         string
	Weapon       UpgradeWeapon
	Keywords     []Keyword
	Exhaust      bool
}

// UpgradeRestrictions wraps the restrictions object
type UpgradeRestrictions struct {
	Name string
	LDF  string `json:"ldf"`
}

// UpgradeWeapon is a special kind of Weapon
type UpgradeWeapon struct {
	Weapon
	Keywords []Keyword
}

// Keyword is a wrapper for the keyword object
type Keyword struct {
	Name        string
	Description string
}

// CommandCard is a Legion command card
type CommandCard struct {
	LDF         string `json:"ldf"`
	Name        string
	Pips        int
	Orders      string
	Description string
	Weapon      Weapon
}

// Weapon is a generic weapon type to be used across multiple data types
type Weapon struct {
	Name     string
	Range    Range
	Dice     AttackDice
	Keywords []string
}

// Range is a weapon's range
type Range struct {
	From int
	To   int
}

// AttackDice is the pool of dice available to a Weapon
type AttackDice struct {
	White int
	Black int
	Red   int
}

func LoadLegionData() (legionData LegionData) {
	// read the legion-data.json file
	jsonFile, err := os.Open("./legion-data.json")

	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	// unmarshal the json into LegionData
	bytes, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(bytes, &legionData)

	return
}
