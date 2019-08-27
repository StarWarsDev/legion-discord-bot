package commands

import (
	"github.com/StarWarsDev/legion-discord-bot/output"
	"github.com/bwmarrin/discordgo"
)

// Help handles the !help command
func Help() discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "!unit <unit card name>",
			Value: "Displays information about the specified unit",
		},
		{
			Name:  "!upgrade <upgrade card name>",
			Value: "Displays information about the specified upgrade",
		},
		{
			Name:  "!command <command card name>",
			Value: "Displays information about the specified command card",
		},
		{
			Name:  "!search <search term>",
			Value: "Displays search results across all data",
		},
		{
			Name:  "!gonk",
			Value: ":robot:",
		},
		{
			Name:  "!lumpy",
			Value: ":heart:",
		},
		{
			Name:  "!help",
			Value: "This help message",
		},
	}

	info := output.Info("", "")

	info.Fields = fields

	return info
}