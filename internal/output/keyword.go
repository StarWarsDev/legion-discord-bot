package output

import (
	"github.com/StarWarsDev/legion-discord-bot/internal/data"
	"github.com/bwmarrin/discordgo"
)

func Keyword(keyword *data.Keyword) discordgo.MessageEmbed {
	response := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Keyword",
		},
		Title:       keyword.Name,
		Description: keyword.Description,
		Color:       0xF242F5,
	}

	return response
}
