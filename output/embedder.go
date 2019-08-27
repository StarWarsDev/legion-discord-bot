package output

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/StarWarsDev/legion-discord-bot/data"
	"github.com/StarWarsDev/legion-discord-bot/utils"
	"github.com/bwmarrin/discordgo"
)

const (
	colorSuccess     = 0x00AE86
	colorError       = 0xE84A4A
	colorInfo        = 0xF2E82B
	colorCommandCard = 0x4287f5
	colorUpgrade     = 0x609c30
	colorUnit        = 0xffffff
	checkMark        = "✓"
)

// Error returns an Embedder with a level of Success
func Error(title, description string) discordgo.MessageEmbed {
	return newEmbedder(colorError, title, description)
}

// Info returns an Embedder with a level of Success
func Info(title, description string) discordgo.MessageEmbed {
	return newEmbedder(colorInfo, title, description)
}

// CommandCard builds an Embedder for a command card
func CommandCard(card *data.CommandCard) discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField

	pipsField := field("Pips", int2Str(card.Pips))
	ordersField := field("Orders", card.Orders)

	fields = append(fields, &pipsField, &ordersField)

	if diceCount(&card.Weapon.Dice) > 0 {
		weapon := weaponString(&card.Weapon)
		weaponField := field("Weapon", weapon)
		fields = append(fields, &weaponField)
	}

	imageURL := "http://legion-hq.com/commands/" + url.PathEscape(card.Name) + ".png"

	return discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{Name: "Command"},
		Title:       card.Name,
		Description: card.Description,
		Fields:      fields,
		Color:       colorCommandCard,
		Image: &discordgo.MessageEmbedImage{
			URL: imageURL,
		},
	}
}

// Upgrade builds an embedder for an upgrade card
func Upgrade(card *data.Upgrade) discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField

	pointsField := field("Points", int2Str(card.Points))
	slotField := field("Slot", card.Slot)

	fields = append(fields, &pointsField, &slotField)

	if card.Exhaust {
		exhaustibleField := field("Exhaustible", checkMark)
		fields = append(fields, &exhaustibleField)
	}

	if card.Restrictions.Name != "" {
		restrictionsField := field("Restrictions", card.Restrictions.Name)
		fields = append(fields, &restrictionsField)
	}

	if len(card.Keywords) > 0 {
		keywordsField := field("Keywords", joinKeywords(card.Keywords))
		fields = append(fields, &keywordsField)
	}

	if diceCount(&card.Weapon.Dice) > 0 {
		weapon := upgradeWeaponString(&card.Weapon)
		weaponField := field("Weapon", weapon)
		fields = append(fields, &weaponField)
	}

	imageURL := "http://legion-hq.com/upgrades/" + url.PathEscape(card.Name) + ".png"

	if card.LDF == "e11dfocusedfireconfig" {
		imageURL = "http://legion-hq.com/upgrades/E-11D%20Grenade%20Launcher%20Config.png"
	}

	return discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{Name: "Upgrade"},
		Title:       card.Name,
		Description: card.Description,
		Color:       colorUpgrade,
		Fields:      fields,
		Image: &discordgo.MessageEmbedImage{
			URL: imageURL,
		},
	}
}

// Unit builds an embedder for a unit card
func Unit(card *data.Unit) discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField

	pointsField := field("Points", int2Str(card.Points))
	typeField := field("Type", card.Type)
	rankField := field("Rank", card.Rank)
	minisField := field("Minis", int2Str(card.Minis))
	woundsField := field("Wounds", int2Str(card.Wounds))

	fields = append(fields,
		&pointsField,
		&typeField,
		&rankField,
		&minisField,
		&woundsField,
	)

	if card.Resilience > 0 {
		resilienceField := field("Resilience", int2Str(card.Resilience))
		fields = append(fields, &resilienceField)
	} else {
		courageField := field("Courage", int2Str(card.Courage))
		fields = append(fields, &courageField)
	}

	defenseField := field("Defense", card.Defense)
	speedField := field("Speed", int2Str(card.Speed))
	slotsField := field("Slots", strings.Join(card.Slots, ", "))
	keywordsField := field("Keywords", joinKeywords(card.Keywords))
	fields = append(fields,
		&defenseField,
		&speedField,
		&slotsField,
		&keywordsField,
	)

	if len(card.CommandCards) > 0 {
		commandCardsField := field("Command Cards", strings.Join(card.CommandCards, ", "))
		fields = append(fields, &commandCardsField)
	}

	name := card.Name
	if card.Unique {
		name = "• " + name
		uniqueField := field("Unique", checkMark)
		fields = append(fields, &uniqueField)
	}

	if len(card.Weapons) > 0 {
		var weapons []string
		for _, weapon := range card.Weapons {
			weapons = append(weapons, weaponString(&weapon))
		}
		weaponsField := field("Weapons", strings.Join(weapons, "\n\n"))
		fields = append(fields, &weaponsField)
	}

	surgeStr := ""
	if card.Surge.Attack != "" {
		surgeStr = "Attack: " + card.Surge.Attack
	}

	if card.Surge.Defense != "" {
		if surgeStr != "" {
			surgeStr = surgeStr + "\n"
		}
		surgeStr = surgeStr + "Defense: " + card.Surge.Defense
	}

	if surgeStr != "" {
		surgeField := field("Surge", surgeStr)
		fields = append(fields, &surgeField)
	}

	imageURL := "http://legion-hq.com/units/" + url.PathEscape(card.Name) + ".png"
	if card.LDF == "tx225gavwoccupiercombatassaulttank" {
		imageURL = "http://legion-hq.com/units/Assault%20Tank.png"
	}

	return discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{Name: "Unit"},
		Title:       name,
		Description: card.Subtitle,
		Color:       colorUnit,
		Fields:      fields,
		Image:       &discordgo.MessageEmbedImage{URL: imageURL},
	}
}

func field(name, value string) discordgo.MessageEmbedField {
	return discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: true,
	}
}

func newEmbedder(color int, title, description string) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
	}
}

func diceCount(dice *data.AttackDice) int {
	white := dice.White
	black := dice.Black
	red := dice.Red

	return white + black + red
}

func weaponString(weapon *data.Weapon) string {
	var weaponInfo []string
	if len(weapon.Name) > 0 {
		weaponInfo = append(weaponInfo, "  "+weapon.Name)
	}

	weaponInfo = append(weaponInfo, utils.WithTemplate("Range: %d - %d", weapon.Range.From, weapon.Range.To))
	weaponInfo = append(weaponInfo, utils.WithTemplate("Dice: %s", utils.DiceString(&weapon.Dice)))

	if len(weapon.Keywords) > 0 {
		keywords := strings.Join(weapon.Keywords, ", ")
		weaponInfo = append(weaponInfo, "Keywords: "+keywords)
	}

	return strings.Join(weaponInfo, "\n")
}

func upgradeWeaponString(weapon *data.UpgradeWeapon) string {
	var weaponInfo []string
	if len(weapon.Name) > 0 {
		weaponInfo = append(weaponInfo, "  "+weapon.Name)
	}

	weaponInfo = append(weaponInfo, utils.WithTemplate("Range: %d - %d", weapon.Range.From, weapon.Range.To))
	weaponInfo = append(weaponInfo, utils.WithTemplate("Dice: %s", utils.DiceString(&weapon.Dice)))

	if len(weapon.Keywords) > 0 {

		weaponInfo = append(weaponInfo, "Keywords: "+joinKeywords(weapon.Keywords))
	}

	return strings.Join(weaponInfo, "\n")
}

func joinKeywords(keywords []data.Keyword) string {
	var ks []string
	for _, keyword := range keywords {
		ks = append(ks, keyword.Name)
	}
	return strings.Join(ks, ", ")
}

func int2Str(v int) string {
	return fmt.Sprintf("%d", v)
}
