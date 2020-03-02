package output

import (
	"github.com/bwmarrin/discordgo"
)

const (
	colorError = 0xE84A4A
	colorInfo  = 0xF2E82B
)

// Error returns an Embedder with a level of Success
func Error(title, description string) discordgo.MessageEmbed {
	return newEmbedder(colorError, title, description)
}

// Info returns an Embedder with a level of Success
func Info(title, description string) discordgo.MessageEmbed {
	return newEmbedder(colorInfo, title, description)
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
