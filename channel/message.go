package channel

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/StarWarsDev/legion-discord-bot/commands"
	"github.com/StarWarsDev/legion-discord-bot/internal/data"
	"github.com/StarWarsDev/legion-discord-bot/output"
	"github.com/StarWarsDev/legion-discord-bot/utils"
	"github.com/bwmarrin/discordgo"
)

func messageSendEmbed(s *discordgo.Session, m *discordgo.MessageCreate, embed *discordgo.MessageEmbed) {
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Embedded item as follows")
		fmt.Printf("%d Fields\n", len(embed.Fields))

		b, err := json.Marshal(embed)
		if err != nil {
			fmt.Println("Could not unmarshal data")
			fmt.Println(err)
		}

		fmt.Println(string(b))

		outputError := output.Error("Failed to render", utils.WithTemplate("There was a problem rendering the %s \"%s\"", embed.Author.Name, embed.Title))
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &outputError)
	}
}

// NewMessageHandler returns a handler function with a bound context for *lookup.Util and *search.Util access
func NewMessageHandler(client *data.ArchivesClient) interface{} {
	// create the message handler function as a variable so we can return it with the parent's context
	handler := func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "!help" {
			helpContent := commands.Help()
			messageSendEmbed(s, m, &helpContent)
		}

		if strings.HasPrefix(m.Content, "!keyword") {
			// split the message so we can interrogate it for intention
			cmdParts := strings.Split(m.Content, " ")
			if len(cmdParts) > 1 {
				// by default, assume the field we want to search is "name"
				field := "name"
				// by default, assume anything after the command is the search term for the name field
				term := strings.TrimSpace(strings.Replace(m.Content, "!keyword", "", 1))
				// if there are more than 2 parts in the message assume that the search is more complex
				if len(cmdParts) > 2 {
					// assume that the field to be searched is following the command
					field = cmdParts[1]
					// assume that the rest of the body is the search term
					term = strings.TrimSpace(strings.Replace(term, field, "", 1))
				}
				// turn all terms into case insensitive regex searches
				term = fmt.Sprintf("(?i)(%s)", term)
				// get the keywords from the archives client
				keywords := client.GetKeywords(field, term)
				// TODO: for each result, send an embedded response message
				fmt.Println(keywords)
			}
		}

		//if strings.HasPrefix(m.Content, "!unit") {
		//	response := commands.Unit(m, lookupUtil)
		//	messageSendEmbed(s, m, &response)
		//}

		//if strings.HasPrefix(m.Content, "!upgrade") {
		//	response := commands.Upgrade(m, lookupUtil)
		//	messageSendEmbed(s, m, &response)
		//}

		//if strings.HasPrefix(m.Content, "!command") {
		//	response := commands.Command(m, lookupUtil)
		//
		//	messageSendEmbed(s, m, &response)
		//}

		//if strings.HasPrefix(m.Content, "!search") {
		//	searchText := strings.Replace(m.Content, "!search", "", 1)
		//	searchText = strings.TrimSpace(searchText)
		//
		//	if strings.ToLower(searchText) == "sexy rexy" {
		//		searchText = "clone captain rex"
		//	}
		//
		//	if strings.ToLower(searchText) == "help" {
		//		outputInfo := output.Info("Search Help", "Find more info on how to structure your search here: http://blevesearch.com/docs/Query-String-Query/")
		//		messageSendEmbed(s, m, &outputInfo)
		//	} else {
		//
		//		if searchText == "" {
		//			response := m.Author.Mention() + ", the `!search` command requires a search term. Please try again using this format `!search <search term>`"
		//			outputError := output.Error("Bad input", response)
		//			messageSendEmbed(s, m, &outputError)
		//		} else {
		//			embeddedResults := searchUtil.FullSearch(searchText)
		//			for _, embed := range embeddedResults {
		//				messageSendEmbed(s, m, &embed)
		//			}
		//		}
		//	}
		//}
	}

	return handler
}
