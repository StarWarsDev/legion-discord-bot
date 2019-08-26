package lookup

import (
	"github.com/StarWarsDev/legion-discord-bot/data"
	"github.com/StarWarsDev/legion-discord-bot/utils"
)

// LookupCommand finds a command card by its name
func (util *Util) LookupCommand(commandName string) data.CommandCard {
	var commandCard data.CommandCard
	cName := utils.CleanName(commandName)

	for _, card := range util.legionData.CommandCards {
		cardName := utils.CleanName(card.Name)
		if cardName == cName {
			commandCard = card
			break
		}
	}

	return commandCard
}

// LookupCommandCardByLdf finds a command card by its LDF value
func (util *Util) LookupCommandCardByLdf(ldf string) data.CommandCard {
	var commandCard data.CommandCard

	for _, card := range util.legionData.CommandCards {
		if ldf == card.LDF {
			commandCard = card
			break
		}
	}

	return commandCard
}
