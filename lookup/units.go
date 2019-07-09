package lookup

import (
	"strings"

	"github.com/StarWarsDev/legion-discord-bot/utils"
)

// LookupUnit finds a unit by name and returns a formatted string slice of related data
func (util *Util) LookupUnit(unitName string) []string {
	cardInfo := []string{}
	uName := utils.CleanName(unitName)

	units := util.legionData.Units.Flattened()
	for _, unit := range units {
		name := utils.CleanName(unit.Name)

		if name == uName {
			cardInfo = []string{
				"Name: " + unit.Name,
				utils.WithTemplate("Points: %d", unit.Points),
				"Type: " + unit.Type,
				"Rank: " + unit.Rank,
				utils.WithTemplate("Minis: %d", unit.Minis),
				utils.WithTemplate("Wounds: %d", unit.Wounds),
			}

			if unit.Type == "Trooper" {
				cardInfo = append(cardInfo, utils.WithTemplate("Courage: %v", utils.Courage(unit.Courage)))
			}

			if unit.Type == "Vehicle" {
				cardInfo = append(cardInfo, utils.WithTemplate("Resilence: %d", unit.Resilience))
			}

			if (unit.Surge.Attack + unit.Surge.Defense) != "" {
				cardInfo = append(cardInfo, utils.WithTemplate("Surge: %s", utils.SurgeString(unit.Surge)))
			}

			cardInfo = append(cardInfo, []string{
				utils.WithTemplate("Speed: %d", unit.Speed),
				"Slots: " + strings.Join(unit.Slots, ", "),
			}...)

			if len(unit.Keywords) > 0 {
				var kwords []string
				for _, keyword := range unit.Keywords {
					kwords = append(kwords, keyword.Name)
				}
				cardInfo = append(cardInfo, "Keywords: "+strings.Join(kwords, ", "))
			}

			if len(unit.CommandCards) > 0 {
				var commandCards []string

				for _, ldf := range unit.CommandCards {
					card := util.LookupCommandCardByLdf(ldf)
					if card != nil {
						commandCards = append(commandCards, card.Name)
					}
				}

				cardInfo = append(
					cardInfo,
					utils.WithTemplate("Command cards: %s", strings.Join(commandCards, ", ")),
				)
			}

			if len(unit.Weapons) > 0 {
				cardInfo = append(cardInfo, "Weapon(s):")
				for _, weapon := range unit.Weapons {
					cardInfo = append(cardInfo, "\n  "+weapon.Name)
					cardInfo = append(cardInfo, utils.WithTemplate(
						"    Range: %d - %d",
						weapon.Range.From,
						weapon.Range.To,
					))
					cardInfo = append(cardInfo, utils.WithTemplate(
						"    Dice: %s",
						utils.DiceString(&weapon.Dice),
					))

					if len(weapon.Keywords) > 0 {
						cardInfo = append(cardInfo, "    Keywords: "+strings.Join(weapon.Keywords, ", "))
					}
				}
			}

			return cardInfo
		}
	}

	if len(cardInfo) == 0 {
		cardInfo = []string{"Nothing found for \"" + unitName + "\""}
	}

	return cardInfo
}
