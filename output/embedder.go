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

	pipsField := Field("Pips", int2Str(card.Pips))
	ordersField := Field("Orders", card.Orders)

	fields = append(fields, &pipsField, &ordersField)

	if diceCount(&card.Weapon.Dice) > 0 {
		weapon := weaponString(&card.Weapon)
		weaponField := Field("Weapon", weapon)
		fields = append(fields, &weaponField)
	}

	imageURL := "http://legion-hq.com/commands/" + url.PathEscape(card.Name) + ".png"

	if card.LDF == "andnowyouwilldie" {
		imageURL = "http://legion-hq.com/commands/And%20Now%20You%20Will%20Die.png"
	}

	if card.LDF == "giveintoyouranger" {
		imageURL = "http://legion-hq.com/commands/Give%20In%20To%20Your%20Anger.png"
	}

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

	pointsField := Field("Points", int2Str(card.Points))
	slotField := Field("Slot", card.Slot)

	fields = append(fields, &pointsField, &slotField)

	if card.Exhaust {
		exhaustibleField := Field("Exhaustible", checkMark)
		fields = append(fields, &exhaustibleField)
	}

	if card.Restrictions.Name != "" {
		restrictionsField := Field("Restrictions", card.Restrictions.Name)
		fields = append(fields, &restrictionsField)
	}

	if len(card.Keywords) > 0 {
		keywordsField := Field("Keywords", joinKeywords(card.Keywords))
		fields = append(fields, &keywordsField)
	}

	if diceCount(&card.Weapon.Dice) > 0 {
		weapon := upgradeWeaponString(&card.Weapon)
		weaponField := Field("Weapon", weapon)
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

	pointsField := Field("Points", int2Str(card.Points))
	typeField := Field("Type", card.Type)
	rankField := Field("Rank", card.Rank)
	minisField := Field("Minis", int2Str(card.Minis))
	woundsField := Field("Wounds", int2Str(card.Wounds))

	fields = append(fields,
		&pointsField,
		&typeField,
		&rankField,
		&minisField,
		&woundsField,
	)

	if card.Resilience > 0 {
		resilienceField := Field("Resilience", int2Str(card.Resilience))
		fields = append(fields, &resilienceField)
	} else {
		courageField := Field("Courage", int2Str(card.Courage))
		fields = append(fields, &courageField)
	}

	defenseField := Field("Defense", card.Defense)
	speedField := Field("Speed", int2Str(card.Speed))
	slotsField := Field("Slots", strings.Join(card.Slots, ", "))
	keywordsField := Field("Keywords", joinKeywords(card.Keywords))
	fields = append(fields,
		&defenseField,
		&speedField,
		&slotsField,
		&keywordsField,
	)

	if len(card.CommandCards) > 0 {
		commandCardsField := Field("Command Cards", strings.Join(card.CommandCards, ", "))
		fields = append(fields, &commandCardsField)
	}

	name := card.Name
	if card.Unique {
		name = "• " + name
		uniqueField := Field("Unique", checkMark)
		fields = append(fields, &uniqueField)
	}

	if len(card.Weapons) > 0 {
		var weapons []string
		for _, weapon := range card.Weapons {
			weapons = append(weapons, weaponString(&weapon))
		}
		weaponsField := Field("Weapons", strings.Join(weapons, "\n\n"))
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
		surgeField := Field("Surge", surgeStr)
		fields = append(fields, &surgeField)
	}

	imageURL := "http://legion-hq.com/units/" + url.PathEscape(card.Name) + ".png"
	if card.LDF == "tx225gavwoccupiercombatassaulttank" {
		imageURL = "http://legion-hq.com/units/Assault%20Tank.png"
	}
	if card.LDF == "rebelofficer" {
		imageURL = "http://legion-hq.com/units/Rebel%20Commander.png"
	}
	if card.LDF == "imperialofficer" {
		imageURL = "http://legion-hq.com/units/Imperial%20Commander.png"
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

// Field creates an embedded field
func Field(name, value string) discordgo.MessageEmbedField {
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
