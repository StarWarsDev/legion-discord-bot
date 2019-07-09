package lookup

import (
	"strings"

	"github.com/StarWarsDev/legion-discord-bot/utils"
)

// LookupUpgrade finds an upgrade based on its name
func (util *Util) LookupUpgrade(upgradeName string) []string {
	cardInfo := []string{}
	upgrName := utils.CleanName(upgradeName)

	upgrades := util.legionData.Upgrades.Flattened()
	for _, upgrade := range upgrades {
		uName := utils.CleanName(upgrade.Name)

		if upgrName == uName {
			cardInfo = []string{
				"Name: " + upgrade.Name,
				utils.WithTemplate("Points: %d", upgrade.Points),
				"Slot: " + upgrade.Slot,
			}

			if upgrade.Exhaust {
				cardInfo = append(cardInfo, "Exhaustable")
			}

			if len(upgrade.Description) > 0 {
				cardInfo = append(cardInfo, "Description: "+upgrade.Description)
			}

			if upgrade.Restrictions.Name != "" {
				cardInfo = append(cardInfo, "Restrictions: "+upgrade.Restrictions.Name)
			}

			if len(upgrade.Keywords) > 0 {
				keywords := []string{}
				for _, keyword := range upgrade.Keywords {
					keywords = append(keywords, keyword.Name)
				}
				cardInfo = append(cardInfo, "Keywords: "+strings.Join(keywords, ","))
			}

			diceCount := upgrade.Weapon.Dice.White + upgrade.Weapon.Dice.Black + upgrade.Weapon.Dice.Red
			if diceCount > 0 {
				weaponInfo := []string{"Weapon:"}
				if len(upgrade.Weapon.Name) > 0 {
					weaponInfo = append(weaponInfo, "  "+upgrade.Weapon.Name)
				}

				weaponInfo = append(weaponInfo, utils.WithTemplate("  Range: %d - %d", upgrade.Weapon.Range.From, upgrade.Weapon.Range.To))
				weaponInfo = append(weaponInfo, utils.WithTemplate("  Dice: %s", utils.DiceString(&upgrade.Weapon.Dice)))

				if len(upgrade.Weapon.Keywords) > 0 {
					keywords := []string{}
					for _, keyword := range upgrade.Weapon.Keywords {
						keywords = append(keywords, keyword.Name)
					}
					weaponInfo = append(weaponInfo, "  Keywords: "+strings.Join(keywords, ","))
				}

				cardInfo = append(cardInfo, weaponInfo...)
			}

			return cardInfo
		}
	}

	if len(cardInfo) == 0 {
		cardInfo = []string{"Nothing found for \"" + upgradeName + "\""}
	}

	return cardInfo
}
