package commands

import (
	"github.com/StarWarsDev/legion-discord-bot/lookup"
	"github.com/StarWarsDev/legion-discord-bot/output"
	"github.com/bwmarrin/discordgo"
	"strings"
)

// Unit handles the !unit command
func Unit(m *discordgo.MessageCreate, lookupUtil *lookup.Util) discordgo.MessageEmbed {
	unitName := strings.Replace(m.Content, "!unit", "", 1)
	unitName = strings.TrimSpace(unitName)
	var response discordgo.MessageEmbed
	if len(unitName) == 0 {
		response = output.Error(
			"Bad input",
			m.Author.Mention()+", the `!unit` command requires a unit card name. Please try again using this format `!unit <unit card name>`",
		)
	} else {
		if strings.ToLower(unitName) == "sexy rexy" {
			unitName = "clone captain rex"
		}

		unit := lookupUtil.LookupUnit(unitName)
		if unit.LDF != "" {
			// replace command card ldf values with names
			if len(unit.CommandCards) > 0 {
				var commandCards []string
				for _, ldf := range unit.CommandCards {
					card := lookupUtil.LookupCommandCardByLdf(ldf)
					if card.LDF != "" {
						commandCards = append(commandCards, card.Name)
					}
				}
				unit.CommandCards = commandCards
			}
			response = output.Unit(&unit)
		} else {
			response = output.Error("No results found", "Nothing found for \""+unitName+"\"")
		}
	}
	return response
}
