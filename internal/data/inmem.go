package data

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	legiondata "github.com/StarWarsDev/go-legion-data"
	legionhq "github.com/StarWarsDev/go-legion-hq"
)

func ImagePathToURL(cardType, imageName string) string {
	return fmt.Sprintf("https://d2b46bduclcqna.cloudfront.net/%sCards/%s", cardType, imageName)
}

type InMemoryClient struct {
	legionHQData legionhq.Data
	extData      *legiondata.Data
}

func (c *InMemoryClient) FetchAllData() {
	log.Println("fetching legion-hq data")
	legionHQData, err := legionhq.GetData()
	if err != nil {
		log.Printf("failed to fetch legion-hq data: %v\n", err)
	} else {
		c.legionHQData = legionHQData
		log.Println("legion-hq data is set")
	}

	// log.Println("fetching legiondata data")
	// extData, err := legiondata.GetData()
	// if err != nil {
	// 	log.Printf("failed to fetch legiondata data: %v\n", err)
	// } else {
	// 	c.extData = extData
	// 	log.Println("legiondata data is set")
	// }
}

// CommandCards returns all data with cardType: command where the field and term matches
func (c *InMemoryClient) CommandCards(field, term string) ([]CommandCard, error) {
	var commandCards []CommandCard
	for _, card := range c.legionHQData.CommandCards() {
		isMatch := false
		switch field {
		case "id":
			isMatch = isExactMatch(card.ID, term)
		case "commander":
			match, err := isRegexpMatch(card.Commander, term)
			if err != nil {
				return nil, err
			}
			isMatch = match
		case "name":
			match, err := isRegexpMatch(card.CardName, term)
			if err != nil {
				return nil, err
			}
			isMatch = match
		case "pips":
			pips, err := strconv.Atoi(term)
			if err != nil {
				return nil, err
			}
			cardPips, err := strconv.Atoi(card.CardSubType)
			if err != nil {
				cardPips = 0
			}

			isMatch = cardPips == pips
		case "faction":
			match, err := isRegexpMatch(card.Faction, term)
			if err != nil {
				return nil, err
			}
			isMatch = match
		default:
			return nil, fmt.Errorf("bad query: field [%s] is not searchable", field)
		}

		if isMatch {
			extCard, _ := c.ExtCommandCard(card.CardName)
			commandCard := CardToCommand(&card, extCard)
			commandCards = append(commandCards, commandCard)
		}
	}
	return commandCards, nil
}

// UnitCards returns all data with cardType: unit where the field and term matches
func (c *InMemoryClient) UnitCards(field, term string) ([]Unit, error) {
	var units []Unit

	for _, card := range c.legionHQData.UnitCards() {
		isMatch := false
		switch field {
		case "name":
			match, err := isRegexpMatch(card.CardName, term)
			if err != nil {
				return nil, err
			}
			isMatch = match
		case "cardType":
			match, err := isRegexpMatch(card.CardType, term)
			if err != nil {
				return nil, err
			}
			isMatch = match
		case "cardSubType":
			match, err := isRegexpMatch(card.CardSubType, term)
			if err != nil {
				return nil, err
			}
			isMatch = match
		case "unique":
			unique, err := strconv.ParseBool(term)
			if err != nil {
				return nil, err
			}
			isMatch = card.IsUnique == unique
		case "requirements":
			for _, requirement := range card.Requirements {
				if !isMatch {
					match, err := isRegexpMatch(requirement, term)
					if err != nil {
						return nil, err
					}
					isMatch = match
				}
			}
		case "rank":
			match, err := isRegexpMatch(card.Rank, term)
			if err != nil {
				return nil, err
			}
			isMatch = match
		case "faction":
			match, err := isRegexpMatch(card.Faction, term)
			if err != nil {
				return nil, err
			}
			isMatch = match
		case "slots":
			for _, slot := range card.UpgradeBar {
				if !isMatch {
					match, err := isRegexpMatch(slot, term)
					if err != nil {
						return nil, err
					}
					isMatch = match
				}
			}
		case "keywords":
			for _, keyword := range card.Keywords {
				if !isMatch {
					match, err := isRegexpMatch(keyword, term)
					if err != nil {
						return nil, err
					}
					isMatch = match
				}
			}
		default:
			return nil, fmt.Errorf("bad query: field [%s] is not searchable", field)
		}

		if isMatch {
			commandCards, err := c.CommandCards("commander", card.CardName)
			if err != nil {
				return nil, err
			}
			extUnit, _ := c.ExtUnit(card.CardName, card.Rank)
			unit := CardToUnit(&card, extUnit)
			if len(commandCards) > 0 {
				unit.CommandCards = commandCards
			}
			units = append(units, unit)
		}
	}

	return units, nil
}

// UpgradeCards returns all data with cardType: upgrade
func (c *InMemoryClient) UpgradeCards(field, term string) ([]Upgrade, error) {
	log.Printf("searching %s for %s", field, term)
	var upgrades []Upgrade

	for _, card := range c.legionHQData.UpgradeCards() {
		isMatch := false
		switch field {
		case "id":
			isMatch = isExactMatch(card.ID, term)
		case "name":
			match, err := isRegexpMatch(card.CardName, term)
			if err != nil {
				return nil, err
			}
			isMatch = match
		case "cardType":
			match, err := isRegexpMatch(card.CardType, term)
			if err != nil {
				return nil, err
			}
			isMatch = match
		case "cardSubType":
			match, err := isRegexpMatch(card.CardSubType, term)
			if err != nil {
				return nil, err
			}
			isMatch = match
		case "unique":
			unique, err := strconv.ParseBool(term)
			if err != nil {
				return nil, err
			}
			isMatch = card.IsUnique == unique
		case "requirements":
			for _, requirement := range card.Requirements {
				if !isMatch {
					match, err := isRegexpMatch(requirement, term)
					if err != nil {
						return nil, err
					}
					isMatch = match
				}
			}
		case "keywords":
			for _, keyword := range card.Keywords {
				if !isMatch {
					match, err := isRegexpMatch(keyword, term)
					if err != nil {
						return nil, err
					}
					isMatch = match
				}
			}
		default:
			return nil, fmt.Errorf("bad query: field [%s] is not searchable", field)
		}
		if isMatch {
			extUpgrade, _ := c.ExtUpgrade(card.CardName)
			upgrade := CardToUpgrade(&card, extUpgrade)
			upgrades = append(upgrades, upgrade)
		}
	}

	return upgrades, nil
}

// Keywords returns a slice of keywords
func (c *InMemoryClient) GetKeywords(nameIn string) ([]Keyword, error) {
	var keywords []Keyword

	for _, collection := range []map[string]string{c.legionHQData.KeywordDict, c.legionHQData.AdditionalKeywords} {
		for name, description := range collection {
			match, _ := isRegexpMatch(name, nameIn)
			if match {
				keyword := Keyword{
					Name:        name,
					Description: description,
				}
				keywords = append(keywords, keyword)
			}
		}
	}

	return keywords, nil
}

func isExactMatch(subject, term string) bool {
	return subject == term
}

func isRegexpMatch(subject, pattern string) (bool, error) {
	return regexp.MatchString(pattern, subject)
}

func (c *InMemoryClient) ExtUnit(name, rank string) (*legiondata.Unit, error) {
	if c.extData == nil {
		return nil, nil
	}

	for _, unit := range c.extData.Units {
		if strings.ToLower(unit.Name) == strings.ToLower(name) &&
			strings.ToLower(unit.Rank) == strings.ToLower(rank) {
			return &unit, nil
		}
	}

	return nil, fmt.Errorf("cound not find unit named %s", name)
}

func (c *InMemoryClient) ExtUpgrade(name string) (*legiondata.Upgrade, error) {
	if c.extData == nil {
		return nil, nil
	}
	for _, upgrade := range c.extData.Upgrades {
		if upgrade.Name == name {
			return &upgrade, nil
		}
	}

	return nil, fmt.Errorf("cound not find upgrade named %s", name)
}

func (c *InMemoryClient) ExtCommandCard(name string) (*legiondata.CommandCard, error) {
	if c.extData == nil {
		return nil, nil
	}
	for _, card := range c.extData.CommandCards {
		if card.Name == name {
			return &card, nil
		}
	}

	return nil, fmt.Errorf("cound not find command card named %s", name)
}

func CardToCommand(card *legionhq.Card, extCard *legiondata.CommandCard) CommandCard {
	pips, err := strconv.Atoi(card.CardSubType)
	if err != nil {
		pips = 0
	}

	commandCard := CommandCard{
		Name:      card.CardName,
		Image:     ImagePathToURL(card.CardType, card.ImageLocation),
		Commander: card.Commander,
		Faction:   card.Faction,
		Pips:      pips,
	}

	// loop through the requirements so we can filter out blanks
	for _, requirement := range card.Requirements {
		if requirement != "" {
			commandCard.Requirements = append(commandCard.Requirements, requirement)
		}
	}

	// loop through the keywords so we can filter out blanks
	for _, keyword := range card.Keywords {
		if keyword != "" {
			commandCard.Keywords = append(commandCard.Keywords, keyword)
		}
	}

	if extCard != nil {
		commandCard.Orders = extCard.Orders
		commandCard.Text = extCard.Description
		if extCard.Weapon != nil {
			commandCard.Weapon = &Weapon{
				Name: extCard.Weapon.Name,
				Range: Range{
					From: extCard.Weapon.Range.From,
					To:   extCard.Weapon.Range.To,
				},
				Keywords: extCard.Weapon.Keywords,
				Dice: Dice{
					Black: extCard.Weapon.Dice.Black,
					Red:   extCard.Weapon.Dice.Red,
					White: extCard.Weapon.Dice.White,
				},
			}

			if extCard.Weapon.Surge != nil {
				commandCard.Weapon.Surge = &Surge{
					Attack:  extCard.Weapon.Surge.Attack,
					Defense: extCard.Weapon.Surge.Defense,
				}
			}
		}
	}

	return commandCard
}

// CardToUnit converts a legionhq card into a Unit
func CardToUnit(card *legionhq.Card, extUnit *legiondata.Unit) Unit {
	unit := Unit{
		Name:    card.CardName,
		Type:    card.CardSubType,
		Image:   ImagePathToURL(card.CardType, card.ImageLocation),
		Unique:  card.IsUnique,
		Cost:    card.Cost,
		Rank:    card.Rank,
		Faction: card.Faction,
	}

	// loop through the slots to filter out blanks
	for _, slot := range card.UpgradeBar {
		if slot != "" {
			unit.Slots = append(unit.Slots, slot)
		}
	}

	// loop through the requirements so we can filter out blanks
	for _, requirement := range card.Requirements {
		if requirement != "" {
			unit.Requirements = append(unit.Requirements, requirement)
		}
	}

	// loop through the keywords so we can filter out blanks
	for _, keyword := range card.Keywords {
		if keyword != "" {
			unit.Keywords = append(unit.Keywords, keyword)
		}
	}

	if extUnit != nil {
		// enrich unit data
		unit.Wounds = extUnit.Wounds
		unit.Courage = extUnit.Courage
		unit.Resilience = extUnit.Resilience
		unit.Surge = &Surge{
			Attack:  extUnit.Surge.Attack,
			Defense: extUnit.Surge.Defense,
		}
		unit.Entourage = extUnit.Entourage
		for _, weapon := range extUnit.Weapons {
			weap := Weapon{
				Name: weapon.Name,
				Range: Range{
					From: weapon.Range.From,
					To:   weapon.Range.To,
				},
				Keywords: weapon.Keywords,
				Dice: Dice{
					Black: weapon.Dice.Black,
					Red:   weapon.Dice.Red,
					White: weapon.Dice.White,
				},
			}

			if weapon.Surge != nil {
				weap.Surge = &Surge{
					Attack:  weapon.Surge.Attack,
					Defense: weapon.Surge.Defense,
				}
			}

			unit.Weapons = append(unit.Weapons, weap)
		}
	}

	return unit
}

func CardToUpgrade(card *legionhq.Card, extUpgrade *legiondata.Upgrade) Upgrade {
	upgrade := Upgrade{
		Type:   card.CardSubType,
		Name:   card.CardName,
		Image:  ImagePathToURL(card.CardType, card.ImageLocation),
		Unique: card.IsUnique,
		Cost:   card.Cost,
	}

	// loop through the requirements so we can filter out blanks
	for _, requirement := range card.Requirements {
		if requirement != "" {
			upgrade.Requirements = append(upgrade.Requirements, requirement)
		}
	}

	// loop through the keywords so we can filter out blanks
	for _, keyword := range card.Keywords {
		if keyword != "" {
			upgrade.Keywords = append(upgrade.Keywords, keyword)
		}
	}

	if extUpgrade != nil {
		upgrade.Exhaust = extUpgrade.Exhaust != nil && *extUpgrade.Exhaust
		upgrade.UnitTypeExclusions = extUpgrade.UnitTypeExclusions
		upgrade.Text = extUpgrade.Description
		if extUpgrade.Weapon != nil {
			upgrade.Weapon = &Weapon{
				Name: extUpgrade.Weapon.Name,
				Range: Range{
					From: extUpgrade.Weapon.Range.From,
					To:   extUpgrade.Weapon.Range.To,
				},
				Keywords: extUpgrade.Weapon.Keywords,
				Dice: Dice{
					Black: extUpgrade.Weapon.Dice.Black,
					Red:   extUpgrade.Weapon.Dice.Red,
					White: extUpgrade.Weapon.Dice.White,
				},
			}

			if extUpgrade.Weapon.Surge != nil {
				upgrade.Weapon.Surge = &Surge{
					Attack:  extUpgrade.Weapon.Surge.Attack,
					Defense: extUpgrade.Weapon.Surge.Defense,
				}
			}
		}
	}

	return upgrade
}
