package output

import (
	"fmt"
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
	maxTitle         = 256
	maxValue         = 1024
	maxDescription   = 2048
	maxEmbedLength   = 6000
	maxNumFields     = 25
	checkMark        = "✓"
)

// Success returns an Embedder with a level of Success
func Success(title, description string) *discordgo.MessageEmbed {
	return newEmbedder(colorSuccess, title, description)
}

// Error returns an Embedder with a level of Success
func Error(title, description string) *discordgo.MessageEmbed {
	return newEmbedder(colorError, title, description)
}

// Info returns an Embedder with a level of Success
func Info(title, description string) *discordgo.MessageEmbed {
	return newEmbedder(colorInfo, title, description)
}

// CommandCard builds an Embedder for a command card
func CommandCard(card *data.CommandCard) *discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{
		field("Pips", int2Str(card.Pips)),
		field("Orders", card.Orders),
	}

	if diceCount(&card.Weapon.Dice) > 0 {
		weapon := weaponString(&card.Weapon)
		fields = append(fields, field("Weapon", weapon))
	}

	return &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{Name: "Command"},
		Title:       card.Name,
		Description: card.Description,
		Fields:      fields,
		Color:       colorCommandCard,
	}
}

// Upgrade builds an embedder for an upgrade card
func Upgrade(card *data.Upgrade) *discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{
		field("Points", int2Str(card.Points)),
		field("Slot", card.Slot),
	}

	if card.Exhaust {
		fields = append(fields, field("Exhaustable", checkMark))
	}

	if card.Restrictions.Name != "" {
		fields = append(fields, field("Restrictions", card.Restrictions.Name))
	}

	if len(card.Keywords) > 0 {
		fields = append(fields, field("Keywords", joinKeywords(card.Keywords)))
	}

	if diceCount(&card.Weapon.Dice) > 0 {
		weapon := upgradeWeaponString(&card.Weapon)
		fields = append(fields, field("Weapon", weapon))
	}

	return &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{Name: "Upgrade"},
		Title:       card.Name,
		Description: card.Description,
		Color:       colorUpgrade,
		Fields:      fields,
	}
}

// Unit builds an embedder for a unit card
func Unit(card *data.Unit) *discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{
		field("Points", int2Str(card.Points)),
		field("Type", card.Type),
		field("Rank", card.Rank),
		field("Minis", int2Str(card.Minis)),
		field("Wounds", int2Str(card.Wounds)),
	}

	if card.Resilience > 0 {
		fields = append(fields, field("Resilience", int2Str(card.Resilience)))
	} else {
		fields = append(fields, field("Courage", int2Str(card.Courage)))
	}

	fields = append(fields,
		field("Defense", card.Defense),
		field("Speed", int2Str(card.Speed)),
		field("Slots", strings.Join(card.Slots, ", ")),
		field("Keywords", joinKeywords(card.Keywords)),
		field("Command Cards", strings.Join(card.CommandCards, ", ")),
	)

	name := card.Name
	if card.Unique {
		name = "• " + name
		fields = append(fields, field("Unique", checkMark))
	}

	if len(card.Weapons) > 0 {
		weapons := []string{}
		for _, weapon := range card.Weapons {
			weapons = append(weapons, weaponString(weapon))
		}
		fields = append(fields, field("Weapons", strings.Join(weapons, "\n\n")))
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
		fields = append(fields, field("Surge", surgeStr))
	}

	return &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{Name: "Unit"},
		Title:       name,
		Description: card.Subtitle,
		Color:       colorUnit,
		Fields:      fields,
	}
}

func field(name, value string) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: true,
	}
}

func newEmbedder(color int, title, description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
	}
}

func diceCount(dice *data.AttackDice) int {
	white := dice.White
	black := dice.Black
	red := dice.Red

	return (white + black + red)
}

func weaponString(weapon *data.Weapon) string {
	weaponInfo := []string{}
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
	weaponInfo := []string{}
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

func joinKeywords(keywords []*data.Keyword) string {
	ks := []string{}
	for _, keyword := range keywords {
		ks = append(ks, keyword.Name)
	}
	return strings.Join(ks, ", ")
}

func int2Str(v int) string {
	return fmt.Sprintf("%d", v)
}
