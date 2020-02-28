package channel

import (
	"encoding/json"
	"fmt"
	"github.com/StarWarsDev/legion-discord-bot/commands"
	"log"
	"strconv"
	"strings"

	"github.com/StarWarsDev/legion-discord-bot/internal/data"
	"github.com/StarWarsDev/legion-discord-bot/output"
	"github.com/StarWarsDev/legion-discord-bot/utils"
	"github.com/bwmarrin/discordgo"
)

func messageSendEmbed(s *discordgo.Session, m *discordgo.MessageCreate, embed *discordgo.MessageEmbed) {
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		log.Println(err)
		log.Println("Embedded item as follows")
		log.Printf("%d Fields\n", len(embed.Fields))

		b, err := json.Marshal(embed)
		if err != nil {
			log.Println("Could not unmarshal data")
			log.Println(err)
		}

		log.Println(string(b))

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

		// split the message so we can interrogate it for intention
		cmdParts := strings.Split(m.Content, " ")
		lenParts := len(cmdParts)

		if lenParts > 0 && strings.HasPrefix(cmdParts[0], "!") {
			command := cmdParts[0]

			field := ""
			term := ""

			// if this is not the help command then parse the incoming message
			if command != "!help" {
				// by default, assume the field we want to search is "name"
				field = "name"
				// set the default term equal to the message body minus the command
				term = strings.TrimSpace(strings.Replace(m.Content, command, "", -1))
				args := strings.Split(term, "=")
				lenArgs := len(args)

				// if there are more than 1 args assume that the search is more complex
				if lenArgs > 1 {
					field = args[0]
					term = strings.TrimSpace(args[1])
				}

				// clean up the term and field, just in case
				field = strings.ToLower(strings.TrimSpace(field))
				term = strings.TrimSpace(term)

				// determine if the term should be a regex
				_, err := strconv.Atoi(term)
				if err != nil {
					// if the term could not be converted to an int then treat is as a regex string
					term = fmt.Sprintf("(?i)(%s)", term)
				}
			}

			log.Println(m.Author.Username+" :", command, field, "=", term)

			if command != "!help" && (field == "" || term == "(?i)()") {
				// this is an invalid command execution, respond with the help command.
				log.Println("empty field or term detected, this won't do")
				sendHelp(s, m)
				// stop trying to process the command
				return
			}

			switch command {
			case "!help":
				sendHelp(s, m)
			case "!keyword":
				// get the keywords from the archives client
				keywords := client.GetKeywords(field, term)
				// for each result, send an embedded response message
				for _, keyword := range keywords {
					response := commands.Keyword(&keyword)
					messageSendEmbed(s, m, &response)
				}
			case "!command":
				// get the command cards from the archives client
				commandCards := client.GetCommandCards(field, term)
				// for each result, send an embedded response message
				for _, card := range commandCards {
					response := commands.Command(&card)
					messageSendEmbed(s, m, &response)
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
	}

	return handler
}

func sendHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpContent := commands.Help()
	messageSendEmbed(s, m, &helpContent)
}
