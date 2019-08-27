package channel

import (
	"encoding/json"
	"fmt"
	"github.com/StarWarsDev/legion-discord-bot/commands"
	"github.com/StarWarsDev/legion-discord-bot/lookup"
	"github.com/StarWarsDev/legion-discord-bot/output"
	"github.com/StarWarsDev/legion-discord-bot/search"
	"github.com/StarWarsDev/legion-discord-bot/utils"
	"github.com/bwmarrin/discordgo"
	"strings"
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
		s.ChannelMessageSendEmbed(m.ChannelID, &outputError)
	}
}

func NewMessageHandler(lookupUtil *lookup.Util, searchUtil *search.Util) interface{} {
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

		if m.Content == "!gonk" {
			gonkContent := commands.Gonk()
			messageSendEmbed(s, m, &gonkContent)
		}

		if m.Content == "!lumpy" {
			randomLumpyURL := commands.Lumpy()

			messageSendEmbed(s, m, &discordgo.MessageEmbed{
				Image: &discordgo.MessageEmbedImage{
					URL: randomLumpyURL,
				},
			})
		}

		if strings.HasPrefix(m.Content, "!unit") {
			response := commands.Unit(m, lookupUtil)
			messageSendEmbed(s, m, &response)
		}

		if strings.HasPrefix(m.Content, "!upgrade") {
			response := commands.Upgrade(m, lookupUtil)
			messageSendEmbed(s, m, &response)
		}

		if strings.HasPrefix(m.Content, "!command") {
			response := commands.Command(m, lookupUtil)

			messageSendEmbed(s, m, &response)
		}

		if strings.HasPrefix(m.Content, "!search") {
			searchText := strings.Replace(m.Content, "!search", "", 1)
			searchText = strings.TrimSpace(searchText)

			if strings.ToLower(searchText) == "help" {
				outputInfo := output.Info("Search Help", "Find more info on how to structure your search here: http://blevesearch.com/docs/Query-String-Query/")
				messageSendEmbed(s, m, &outputInfo)
			} else {

				if searchText == "" {
					response := m.Author.Mention() + ", the `!search` command requires a search term. Please try again using this format `!search <search term>`"
					outputError := output.Error("Bad input", response)
					messageSendEmbed(s, m, &outputError)
				} else {
					embeddedResults := searchUtil.FullSearch(searchText)
					for _, embed := range embeddedResults {
						messageSendEmbed(s, m, &embed)
					}
				}
			}
		}
	}

	return handler
}
