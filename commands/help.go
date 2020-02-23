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
			Value: "Displays the specified unit",
		},
		{
			Name:  "!upgrade <upgrade card name>",
			Value: "Displays the specified upgrade",
		},
		{
			Name:  "!command <command card name>",
			Value: "Displays the specified command card",
		},
		{
			Name:  "!keyword <keyword name>",
			Value: "Displays the specified keyword",
		},
		//{
		//	Name:  "!search <search term>",
		//	Value: "Displays search results across all data",
		//},
		{
			Name:  "!help",
			Value: "This help message",
		},
	}

	info := output.Info("", "")

	info.Fields = fields

	return info
}
