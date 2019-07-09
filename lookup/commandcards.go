package lookup

import (
	"fmt"
	"strings"

	"github.com/StarWarsDev/legion-discord-bot/data"
	"github.com/StarWarsDev/legion-discord-bot/utils"
)

// LookupCommand finds a command card by its name
func (util *Util) LookupCommand(commandName string) []string {
	cardInfo := []string{}
	cName := utils.CleanName(commandName)

	for _, card := range util.legionData.CommandCards {
		cardName := utils.CleanName(card.Name)
		if cardName == cName {
			cardInfo = []string{
				"Name: " + card.Name,
				fmt.Sprintf("Pips: %d", card.Pips),
				"Orders: " + card.Orders,
			}

			if len(card.Description) > 0 {
				cardInfo = append(cardInfo, "Description: "+card.Description)
			}

			diceCount := card.Weapon.Dice.White + card.Weapon.Dice.Black + card.Weapon.Dice.Red
			if diceCount > 0 {
				weaponInfo := []string{"Weapon:"}
				if len(card.Weapon.Name) > 0 {
					weaponInfo = append(weaponInfo, "  "+card.Weapon.Name)
				}

				weaponInfo = append(weaponInfo, utils.WithTemplate("  Range: %d - %d", card.Weapon.Range.From, card.Weapon.Range.To))
				weaponInfo = append(weaponInfo, utils.WithTemplate("  Dice: %s", utils.DiceString(&card.Weapon.Dice)))

				if len(card.Weapon.Keywords) > 0 {
					keywords := strings.Join(card.Weapon.Keywords, ", ")
					weaponInfo = append(weaponInfo, "  Keywords: "+keywords)
				}

				cardInfo = append(cardInfo, weaponInfo...)
			}

			return cardInfo
		}
	}

	if len(cardInfo) == 0 {
		cardInfo = []string{"Nothing found for \"" + commandName + "\""}
	}

	return cardInfo
}

// LookupCommandCardByLdf finds a command card by its LDF value
func (util *Util) LookupCommandCardByLdf(ldf string) *data.CommandCard {
	for _, card := range util.legionData.CommandCards {
		if ldf == card.LDF {
			return card
		}
	}

	return nil
}
