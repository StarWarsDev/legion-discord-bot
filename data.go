package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// LegionData is a container for all Legion data
type LegionData struct {
	Units        []*Unit        `json:"units"`
	Upgrades     []*Upgrade     `json:"upgrades"`
	CommandCards []*CommandCard `json:"commandCards"`
}

// Unit is a Legion unit card
type Unit struct{}

// Upgrade is a Legion upgrade card
type Upgrade struct{}

// CommandCard is a Legion command card
type CommandCard struct {
	Name        string `json:"name"`
	Pips        int    `json:"pips"`
	Orders      string `json:"orders"`
	Description string `json:"description,omitempty"`
	Weapon      Weapon `json:"weapon,omitempty"`
}

// Weapon is a generic weapon type to be used across multiple data types
type Weapon struct {
	Name     string     `json:"name"`
	Range    Range      `json:"range"`
	Dice     AttackDice `json:"dice"`
	Keywords []string   `json:"keywords"`
}

// Range is a weapon's range
type Range struct {
	From int `json:"from,omitempty"`
	To   int `json:"to,omitempty"`
}

// AttackDice is the pool of dice available to a Weapon
type AttackDice struct {
	White int `json:"white,omitempty"`
	Black int `json:"black,omitempty"`
	Red   int `json:"red,omitempty"`
}

func loadLegionData() *LegionData {
	// read the legion-data.json file
	jsonFile, err := os.Open("./legion-data.json")

	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	// unmarshal the json into LegionData
	var legionData LegionData
	bytes, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(bytes, &legionData)

	// return a pointer to the data
	return &legionData
}
