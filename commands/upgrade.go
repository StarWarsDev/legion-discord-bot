package commands

import (
	"github.com/StarWarsDev/legion-discord-bot/internal/data"
	"github.com/StarWarsDev/legion-discord-bot/output"
	"github.com/StarWarsDev/legion-discord-bot/utils"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

// Upgrade handles the !upgrade command
func Upgrade(upgrade *data.Upgrade) discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField

	if upgrade.Text != "" {
		textField := output.Field("Text", upgrade.Text)
		fields = append(fields, &textField)
	}

	costField := output.Field("Cost", strconv.Itoa(upgrade.Cost))
	typeField := output.Field("Type", upgrade.Type)
	exhaustField := output.Field("Exhaust", strconv.FormatBool(upgrade.Exhaust))
	fields = append(
		fields,
		&costField,
		&typeField,
		&exhaustField,
	)

	if len(upgrade.Keywords) > 0 {
		keywordsField := output.Field("Keywords", strings.Join(upgrade.Keywords, ", "))
		fields = append(fields, &keywordsField)
	}

	if len(upgrade.Requirements) > 0 {
		requirementsField := output.Field("Requirements", strings.Join(upgrade.Requirements, ", "))
		fields = append(fields, &requirementsField)
	}

	if len(upgrade.UnitTypeExclusions) > 0 {
		exclusionsField := output.Field("Exclusions", strings.Join(upgrade.UnitTypeExclusions, ", "))
		fields = append(fields, &exclusionsField)
	}

	if upgrade.Weapon != nil {
		weaponField := output.Field("Weapon", upgrade.Weapon.String())
		fields = append(fields, &weaponField)
	}

	response := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Upgrade",
		},
		Fields: fields,
		Color:  0x609c30,
		Title:  upgrade.Name,
		Image: &discordgo.MessageEmbedImage{
			URL: utils.FixURL(upgrade.Image),
		},
	}

	return response
}
