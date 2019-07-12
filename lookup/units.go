package lookup

import (
	"github.com/StarWarsDev/legion-discord-bot/data"
	"github.com/StarWarsDev/legion-discord-bot/utils"
)

// LookupUnit finds a unit by name and returns a formatted string slice of related data
func (util *Util) LookupUnit(unitName string) *data.Unit {
	uName := utils.CleanName(unitName)

	units := util.legionData.Units.Flattened()
	for _, unit := range units {
		name := utils.CleanName(unit.Name)

		if name == uName {
			return unit
		}
	}

	return nil
}

// LookupUpgradeByLdf finds a upgrade card by its LDF value
func (util *Util) LookupUnitByLdf(ldf string) *data.Unit {
	for _, card := range util.legionData.Units.Flattened() {
		if ldf == card.LDF {
			return card
		}
	}

	return nil
}
