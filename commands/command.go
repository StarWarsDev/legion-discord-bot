package commands

import (
	"github.com/StarWarsDev/legion-discord-bot/internal/data"
	"github.com/StarWarsDev/legion-discord-bot/output"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

// Command handles the !command command
func Command(card *data.CommandCard) discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField

	pipsField := output.Field("Pips", strconv.Itoa(card.Pips))
	orders := card.Orders
	if orders == "" {
		orders = card.Commander
	}
	ordersField := output.Field("Orders", orders)

	fields = append(fields, &pipsField, &ordersField)

	if card.Commander != "" {
		commanderField := output.Field("Commander", card.Commander)
		fields = append(fields, &commanderField)
	}

	if card.Faction != "" {
		factionField := output.Field("Faction", card.Faction)
		fields = append(fields, &factionField)
	}

	if len(card.Requirements) > 0 {
		requirements := strings.Join(card.Requirements, ", ")
		if requirements != "" {
			requirementsField := output.Field("Requirements", requirements)
			fields = append(fields, &requirementsField)
		}
	}

	if len(card.Keywords) > 0 {
		keywords := strings.Join(card.Keywords, ", ")
		if keywords != "" {
			keywordsField := output.Field("Keywords", keywords)
			fields = append(fields, &keywordsField)
		}
	}

	if card.Weapon != nil {
		weaponField := output.Field("Weapon", card.Weapon.String())
		fields = append(fields, &weaponField)
	}

	response := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Command Card",
		},
		Color:  0x4287f5,
		Title:  card.Name,
		Fields: fields,
		Image: &discordgo.MessageEmbedImage{
			URL: card.Image,
		},
	}

	if card.Text != "" {
		response.Description = card.Text
	}

	return response
}
