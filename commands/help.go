package commands

import (
	"github.com/StarWarsDev/legion-discord-bot/output"
	"github.com/bwmarrin/discordgo"
)

// Help handles the !help command
func Help() discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "!unit",
			Value: "Displays data about Unit cards",
		},
		{
			Name:  "!upgrade",
			Value: "Displays data about Upgrade cards",
		},
		{
			Name:  "!command",
			Value: "Displays data about Command cards",
		},
		{
			Name:  "!keyword",
			Value: "Displays data about Keywords",
		},
		{
			Name:  "!help",
			Value: "This help message",
		},
	}
	info := output.Info("Help", `**Usage:**
`+"`!COMMAND [FIELD_NAME ~] SEARCH_TERM`"+`

`+"`COMMAND` can be any of the available commands below"+`
`+"`FIELD_NAME`"+` is optional and defaults to "name" and can be any top level field on the item being queried.
`+"`SEARCH_TERM`"+` is treated as a regular expression. All terms are considered _case insensitive_.

**Example:** `+"`!keyword inspire`"+` finds all keywords with "inspire" in the name.

**Example 2:** `+"`!keyword description ~ rally`"+` finds all keywords that have the word "rally" in the description.

**DM**: By default the bot will respond in the channel where the command was issued.
If you like to have the bot send you a DM instead just mention it anywhere in the message.

**Commands:**`)

	info.Fields = fields

	return info
}
