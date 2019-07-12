package lookup

import (
	"github.com/StarWarsDev/legion-discord-bot/data"
	"github.com/StarWarsDev/legion-discord-bot/utils"
)

// LookupUpgrade finds an upgrade based on its name
func (util *Util) LookupUpgrade(upgradeName string) *data.Upgrade {
	upgrName := utils.CleanName(upgradeName)

	upgrades := util.legionData.Upgrades.Flattened()
	for _, upgrade := range upgrades {
		uName := utils.CleanName(upgrade.Name)

		if upgrName == uName {
			return upgrade
		}
	}

	return nil
}

// LookupUpgradeByLdf finds a upgrade card by its LDF value
func (util *Util) LookupUpgradeByLdf(ldf string) *data.Upgrade {
	for _, card := range util.legionData.Upgrades.Flattened() {
		if ldf == card.LDF {
			return card
		}
	}

	return nil
}
