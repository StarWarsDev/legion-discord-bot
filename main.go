package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/StarWarsDev/legion-discord-bot/data"
	"github.com/StarWarsDev/legion-discord-bot/lookup"
	"github.com/StarWarsDev/legion-discord-bot/output"
	"github.com/StarWarsDev/legion-discord-bot/search"
	"github.com/StarWarsDev/legion-discord-bot/utils"
	"github.com/bwmarrin/discordgo"
)

var (
	token      string
	legionData *data.LegionData
	lookupUtil *lookup.Util
	searchUtil *search.Util
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()

	if token == "" {
		// no token was assigned via the t flag, try looking in the environment
		token = os.Getenv("DISCORD_TOKEN")
	}

	if token == "" {
		// still no token, panic!
		panic("No discord token provided! Try passing it with the '-t' flag or setting 'DISCORD_TOKEN' in the environment.")
	}
}

func main() {
	fmt.Println("Hello, World! I am the Discord Legion bot!")

	legionData = data.LoadLegionData()
	lookupUtil = lookup.NewUtil(legionData)
	searchUtil = search.NewUtil(legionData, lookupUtil)

	discord, err := discordgo.New("Bot " + token)

	if err != nil {
		panic(err)
	}

	discord.AddHandler(messageCreate)

	err = discord.Open()

	if err != nil {
		panic(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
	os.RemoveAll(search.IndexKey)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!help" {
		fields := []*discordgo.MessageEmbedField{
			{
				Name:  "!unit <unit card name>",
				Value: "Displays information about the specified unit",
			},
			{
				Name:  "!upgrade <upgrade card name>",
				Value: "Displays information about the specified upgrade",
			},
			{
				Name:  "!command <command card name>",
				Value: "Displays information about the specified command card",
			},
			{
				Name:  "!search <search term>",
				Value: "Displays search results across all data",
			},
			{
				Name:  "!gonk",
				Value: ":robot:",
			},
			{
				Name:  "!lumpy",
				Value: ":heart:",
			},
			{
				Name:  "!help",
				Value: "This help message",
			},
		}

		info := output.Info("", "")

		info.Fields = fields
		channelMessageSendEmbed(s, m, info)
	}

	if m.Content == "!gonk" {
		e := output.Error("GONK!", "")
		e.URL = "https://www.starwars.com/databank/gnk-droid"
		e.Image = &discordgo.MessageEmbedImage{
			URL: "https://lumiere-a.akamaihd.net/v1/images/gnk-droid-main-image_f0d89099.jpeg?region=0%2C80%2C1280%2C720",
		}
		channelMessageSendEmbed(s, m, e)
	}

	if m.Content == "!lumpy" {
		urls := []string{
			"https://media.giphy.com/media/PSnTTyw6BdjHy/giphy.gif",
			"https://media.giphy.com/media/10QqGj0eqGOWIw/giphy.gif",
			"https://media.giphy.com/media/FY5dT7KDV2i0o/giphy.gif",
			"https://media.giphy.com/media/hffHBmxUSfHlm/giphy.gif",
		}

		rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
		randomLumpyURL := urls[rand.Intn(len(urls))]

		channelMessageSendEmbed(s, m, &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				URL: randomLumpyURL,
			},
		})
	}

	if strings.HasPrefix(m.Content, "!unit") {
		unitName := strings.Replace(m.Content, "!unit", "", 1)
		unitName = strings.TrimSpace(unitName)

		var response *discordgo.MessageEmbed
		if len(unitName) == 0 {
			response = output.Error(
				"Bad input",
				m.Author.Mention()+", the `!unit` command requires a unit card name. Please try again using this format `!unit <unit card name>`",
			)
		} else {
			unit := lookupUtil.LookupUnit(unitName)
			if unit != nil {
				// replace command card ldf values with names
				if len(unit.CommandCards) > 0 {
					var commandCards []string
					for _, ldf := range unit.CommandCards {
						card := lookupUtil.LookupCommandCardByLdf(ldf)
						if card != nil {
							commandCards = append(commandCards, card.Name)
						}
					}
					unit.CommandCards = commandCards
				}
				response = output.Unit(unit)
			} else {
				response = output.Error("No results found", "Nothing found for \""+unitName+"\"")
			}
		}

		channelMessageSendEmbed(s, m, response)
	}

	if strings.HasPrefix(m.Content, "!upgrade") {
		upgradeName := strings.Replace(m.Content, "!upgrade", "", 1)
		upgradeName = strings.TrimSpace(upgradeName)

		var response *discordgo.MessageEmbed
		if len(upgradeName) == 0 {
			response = output.Error(
				"Bad input",
				m.Author.Mention()+", the `!upgrade` command requires an upgrade card name. Please try again using this format `!upgrade <upgrade card name>`",
			)
		} else {
			upgrade := lookupUtil.LookupUpgrade(upgradeName)
			if upgrade != nil {
				response = output.Upgrade(upgrade)
			} else {
				response = output.Error("No results found", "Nothing found for \""+upgradeName+"\"")
			}
		}

		channelMessageSendEmbed(s, m, response)
	}

	if strings.HasPrefix(m.Content, "!command") {
		commandName := strings.Replace(m.Content, "!command", "", 1)
		commandName = strings.TrimSpace(commandName)

		var response *discordgo.MessageEmbed
		if len(commandName) == 0 {
			response = output.Error(
				"Bad input",
				m.Author.Mention()+", the `!command` command requires a command card name. Please try again using this format `!command <command card name>`",
			)
		} else {
			command := lookupUtil.LookupCommand(commandName)
			if command != nil {
				response = output.CommandCard(command)
			} else {
				response = output.Error("No results found", "Nothing found for \""+commandName+"\"")
			}
		}

		channelMessageSendEmbed(s, m, response)
	}

	if strings.HasPrefix(m.Content, "!search") {
		searchText := strings.Replace(m.Content, "!search", "", 1)
		searchText = strings.TrimSpace(searchText)

		if strings.ToLower(searchText) == "help" {
			channelMessageSendEmbed(s, m, output.Info("Search Help", "Find more info on how to structure your search here: http://blevesearch.com/docs/Query-String-Query/"))
		} else {

			if searchText == "" {
				response := m.Author.Mention() + ", the `!search` command requires a search term. Please try again using this format `!search <search term>`"
				channelMessageSendEmbed(s, m, output.Error("Bad input", response))
			} else {
				embeddedResults := searchUtil.FullSearch(searchText)
				for _, embed := range embeddedResults {
					channelMessageSendEmbed(s, m, embed)
				}
			}
		}
	}

}

func channelMessageSendEmbed(s *discordgo.Session, m *discordgo.MessageCreate, embed *discordgo.MessageEmbed) {
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

		s.ChannelMessageSendEmbed(m.ChannelID, output.Error("Failed to render", utils.WithTemplate("There was a problem rendering the %s \"%s\"", embed.Author.Name, embed.Title)))
	}
}
