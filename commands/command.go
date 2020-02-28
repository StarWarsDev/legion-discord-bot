package commands

import (
	"github.com/StarWarsDev/legion-discord-bot/internal/data"
	"github.com/StarWarsDev/legion-discord-bot/output"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
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

	//if diceCount(&card.Weapon.Dice) > 0 {
	//	weapon := weaponString(&card.Weapon)
	//	weaponField := field("Weapon", weapon)
	//	fields = append(fields, &weaponField)
	//}

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

	log.Println(response)

	return response
}
