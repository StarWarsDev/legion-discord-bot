package commands

import (
	"github.com/StarWarsDev/legion-discord-bot/lookup"
	"github.com/StarWarsDev/legion-discord-bot/output"
	"github.com/bwmarrin/discordgo"
	"strings"
)

// Upgrade handles the !upgrade command
func Upgrade(m *discordgo.MessageCreate, lookupUtil *lookup.Util) discordgo.MessageEmbed {
	upgradeName := strings.Replace(m.Content, "!upgrade", "", 1)
	upgradeName = strings.TrimSpace(upgradeName)
	var response discordgo.MessageEmbed
	if len(upgradeName) == 0 {
		response = output.Error(
			"Bad input",
			m.Author.Mention()+", the `!upgrade` command requires an upgrade card name. Please try again using this format `!upgrade <upgrade card name>`",
		)
	} else {
		upgrade := lookupUtil.LookupUpgrade(upgradeName)
		if upgrade.LDF != "" {
			response = output.Upgrade(&upgrade)
		} else {
			response = output.Error("No results found", "Nothing found for \""+upgradeName+"\"")
		}
	}
	return response
}
