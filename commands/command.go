package commands

import (
	"github.com/StarWarsDev/legion-discord-bot/lookup"
	"github.com/StarWarsDev/legion-discord-bot/output"
	"github.com/bwmarrin/discordgo"
	"strings"
)

// Command handles the !command command
func Command(m *discordgo.MessageCreate, lookupUtil *lookup.Util) discordgo.MessageEmbed {
	commandName := strings.Replace(m.Content, "!command", "", 1)
	commandName = strings.TrimSpace(commandName)
	var response discordgo.MessageEmbed
	if len(commandName) == 0 {
		response = output.Error(
			"Bad input",
			m.Author.Mention()+", the `!command` command requires a command card name. Please try again using this format `!command <command card name>`",
		)
	} else {
		command := lookupUtil.LookupCommand(commandName)
		if command.LDF != "" {
			response = output.CommandCard(&command)
		} else {
			response = output.Error("No results found", "Nothing found for \""+commandName+"\"")
		}
	}
	return response
}
