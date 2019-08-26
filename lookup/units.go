package lookup

import (
	"github.com/StarWarsDev/legion-discord-bot/data"
	"github.com/StarWarsDev/legion-discord-bot/utils"
)

// LookupUnit finds a unit by name and returns a formatted string slice of related data
func (util *Util) LookupUnit(unitName string) data.Unit {
	uName := utils.CleanName(unitName)
	var unit data.Unit

	units := util.legionData.Units.Flattened()
	for _, card := range units {
		name := utils.CleanName(card.Name)

		if name == uName {
			unit = card
			break
		}
	}

	return unit
}

// LookupUpgradeByLdf finds a upgrade card by its LDF value
func (util *Util) LookupUnitByLdf(ldf string) data.Unit {
	var unit data.Unit
	for _, card := range util.legionData.Units.Flattened() {
		if ldf == card.LDF {
			unit = card
			break
		}
	}

	return unit
}
