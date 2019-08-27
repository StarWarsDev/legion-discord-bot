package commands

import (
	"github.com/StarWarsDev/legion-discord-bot/output"
	"github.com/bwmarrin/discordgo"
)

// Gonk handles the !gonk command
func Gonk() discordgo.MessageEmbed {
	e := output.Error("GONK!", "")
	e.URL = "https://www.starwars.com/databank/gnk-droid"
	e.Image = &discordgo.MessageEmbedImage{
		URL: "https://lumiere-a.akamaihd.net/v1/images/gnk-droid-main-image_f0d89099.jpeg?region=0%2C80%2C1280%2C720",
	}
	return e
}
