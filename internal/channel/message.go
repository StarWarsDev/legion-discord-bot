package channel

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/StarWarsDev/legion-discord-bot/internal/data"
	"github.com/StarWarsDev/legion-discord-bot/internal/output"
	"github.com/bwmarrin/discordgo"
)

var validCommands = map[string]bool{
	"!help":    true,
	"!unit":    true,
	"!command": true,
	"!keyword": true,
	"!upgrade": true,
}

func isValidCommand(command string) bool {
	_, ok := validCommands[command]
	return ok
}

func messageSendEmbed(dm bool, s *discordgo.Session, m *discordgo.MessageCreate, embed *discordgo.MessageEmbed) {
	channelID := channelID(m, dm, s)

	_, err := s.ChannelMessageSendEmbed(channelID, embed)
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

		title := "Failed to render"
		description := fmt.Sprintf("There was a problem rendering the %s \"%s\"", embed.Author.Name, embed.Title)
		errorMessageSendEmbed(title, description, dm, s, m)
	}
}

func channelID(m *discordgo.MessageCreate, dm bool, s *discordgo.Session) string {
	channelID := m.ChannelID
	if dm {
		channel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			log.Println("Error creating private channel with user", m.Author.Username, err)
		} else {
			channelID = channel.ID
		}
	}
	return channelID
}

func errorMessageSendEmbed(title, description string, dm bool, s *discordgo.Session, m *discordgo.MessageCreate) {
	channelID := channelID(m, dm, s)
	outputError := output.Error(title, description)
	_, _ = s.ChannelMessageSendEmbed(channelID, &outputError)
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

		content := m.Content
		isMentioned := false
		for _, user := range m.Mentions {
			if user.Mention() == s.State.User.Mention() {
				isMentioned = true
			}
		}
		if isMentioned {
			content = m.ContentWithMentionsReplaced()
			content = strings.ReplaceAll(content, "@"+s.State.User.Username, "")
			content = strings.TrimSpace(content)
		}

		// if the message content starts with an exclamation then it is probably a command
		if strings.HasPrefix(content, "!") {
			command, field, term := parseCommandFieldAndTerm(content)

			log.Println(m.Author.Username, ":", command, field, "~", term)

			if isValidCommand(command) {
				handleCommand(command, field, term, isMentioned, s, m, client)
			}
		}
	}

	return handler
}

func handleCommand(command string, field string, term string, isMentioned bool, s *discordgo.Session, m *discordgo.MessageCreate, client *data.ArchivesClient) {
	var responses []discordgo.MessageEmbed
	processCommand := true

	if command != "!help" && (field == "" || term == "(?i)()") {
		// this is an invalid command execution, respond with the help command.
		log.Println("empty field or term detected, this won't do")
		responses = append(responses, output.Help())
		processCommand = false
	}

	if processCommand {
		switch command {
		case "!help":
			responses = append(responses, output.Help())
		case "!keyword":
			// get the keywords from the archives client
			keywords := client.GetKeywords(field, term)
			// for each result, send an embedded response message
			for _, keyword := range keywords {
				responses = append(responses, output.Keyword(&keyword))
			}
		case "!command":
			// get the command cards from the archives client
			commandCards := client.GetCommandCards(field, term)
			// for each result, send an embedded response message
			for _, card := range commandCards {
				responses = append(responses, output.Command(&card))
			}
		case "!unit":
			// get the unit cards from the archives
			units := client.GetUnits(field, term)
			// for each result, send an embedded response message
			for _, unit := range units {
				responses = append(responses, output.Unit(&unit))
			}
		case "!upgrade":
			// get the upgrade cards from the archives
			upgrades := client.GetUpgrades(field, term)
			// for each result, send an embedded response message
			for _, upgrade := range upgrades {
				responses = append(responses, output.Upgrade(&upgrade))
			}
		}
	}

	if command != "!help" {
		infoResponse := output.Info(
			fmt.Sprintf("Showing %d Results", len(responses)),
			fmt.Sprintf("Found %d results for `%s %s ~ %s`", len(responses), command, field, term),
		)
		responses = append(responses, infoResponse)
	}

	for _, response := range responses {
		messageSendEmbed(isMentioned, s, m, &response)
	}
}

func parseCommandFieldAndTerm(content string) (string, string, string) {
	cmdParts := strings.Split(content, " ")
	command := cmdParts[0]

	field := ""
	term := ""

	// if this is not the help command then parse the incoming message
	if command != "!help" {
		// by default, assume the field we want to search is "name"
		field = "name"
		// set the default term equal to the message body minus the command
		term = strings.TrimSpace(strings.Replace(content, command, "", -1))
		args := strings.Split(term, "~")
		lenArgs := len(args)

		// if there are more than 1 args assume that the search is more complex
		if lenArgs > 1 {
			field = args[0]
			term = strings.TrimSpace(args[1])
		}

		// clean up the term and field, just in case
		field = strings.ToLower(strings.TrimSpace(field))
		term = ParseTerm(term)
	}
	return command, field, term
}

func ParseTerm(term string) string {
	term = strings.TrimSpace(term)
	canBeRegExp := true

	if _, err := strconv.Atoi(term); err == nil {
		canBeRegExp = false
	}

	if canBeRegExp {
		term = fmt.Sprintf("(?i)(%s)", term)
	}

	return term
}

func HandleSlashCommand(command string, field string, term string, s *discordgo.Session, i *discordgo.InteractionCreate, client *data.ArchivesClient, inMemClient *data.InMemoryClient) {
	//term = ParseTerm(term)
	var responses []*discordgo.MessageEmbed
	processCommand := true

	if processCommand {
		// wrap the term in regexp if not searching for command pips
		if command != "command" && field != "pips" {
			term = fmt.Sprintf("(?i)(%s)", term)
		}
		switch command {
		case "keyword":
			// get the keywords from the archives client
			keywords, err := inMemClient.GetKeywords(term)
			if err != nil {
				log.Printf("there was an error getting keyword data for: [%s]: %v", term, err)
			}
			// for each result, send an embedded response message
			for _, keyword := range keywords {
				response := output.Keyword(&keyword)
				responses = append(responses, &response)
			}
		case "command":
			// get the command cards from the archives client
			commandCards, _ := inMemClient.CommandCards(field, term)
			// for each result, send an embedded response message
			for _, card := range commandCards {
				response := output.Command(&card)
				responses = append(responses, &response)
			}
		case "unit":
			// get the unit cards from the archives
			units, _ := inMemClient.UnitCards(field, term)
			// for each result, send an embedded response message
			for _, unit := range units {
				response := output.Unit(&unit)
				responses = append(responses, &response)
			}
		case "upgrade":
			// get the upgrade cards from the archives
			upgrades, _ := inMemClient.UpgradeCards(field, term)
			// for each result, send an embedded response message
			for _, upgrade := range upgrades {
				response := output.Upgrade(&upgrade)
				responses = append(responses, &response)
			}
		}
	}

	infoResponse := output.Info(
		fmt.Sprintf("Showing %d Results", len(responses)),
		fmt.Sprintf("Found %d results for `%s %s ~ %s`", len(responses), command, field, term),
	)
	responses = append(responses, &infoResponse)

	interactionSendEmbed(s, i, responses)
}

func interactionSendEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embeds []*discordgo.MessageEmbed) {
	// channelID := i.ChannelID

	log.Println("sending response back to discord")

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		},
	})

	// _, err := s.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		log.Println(err)
		log.Println("Embedded item as follows")
		for _, embed := range embeds {
			log.Printf("%d Fields\n", len(embed.Fields))

			b, err := json.Marshal(embed)
			if err != nil {
				log.Println("Could not unmarshal data")
				log.Println(err)
			}

			log.Println(string(b))
		}
		// title := "Failed to render"
		// description := fmt.Sprintf("There was a problem rendering the %s \"%s\"", embed.Author.Name, embed.Title)
		// errorMessageSendEmbed(title, description, dm, s, m)
	}
}
