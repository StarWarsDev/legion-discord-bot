package main

import (
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
			&discordgo.MessageEmbedField{
				Name:  "!unit <unit card name>",
				Value: "Displays information about the specified unit",
			},
			&discordgo.MessageEmbedField{
				Name:  "!upgrade <upgrade card name>",
				Value: "Displays information about the specified upgrade",
			},
			&discordgo.MessageEmbedField{
				Name:  "!command <command card name>",
				Value: "Displays information about the specified command card",
			},
			&discordgo.MessageEmbedField{
				Name:  "!search <search term>",
				Value: "Displays search results across all data",
			},
			&discordgo.MessageEmbedField{
				Name:  "!gonk",
				Value: ":robot:",
			},
			&discordgo.MessageEmbedField{
				Name:  "!lumpy",
				Value: ":heart:",
			},
			&discordgo.MessageEmbedField{
				Name:  "!help",
				Value: "This help message",
			},
		}

		info := output.Info("", "")

		info.Fields = fields
		s.ChannelMessageSendEmbed(m.ChannelID, info)
	}

	if m.Content == "!gonk" {
		e := output.Error("GONK!", "")
		e.URL = "https://www.starwars.com/databank/gnk-droid"
		e.Image = &discordgo.MessageEmbedImage{
			URL: "https://lumiere-a.akamaihd.net/v1/images/gnk-droid-main-image_f0d89099.jpeg?region=0%2C80%2C1280%2C720",
		}
		s.ChannelMessageSendEmbed(m.ChannelID, e)
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

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				URL: randomLumpyURL,
			},
		})
	}

	if strings.HasPrefix(m.Content, "!unit") {
		unitName := strings.Replace(m.Content, "!unit", "", 1)
		unitName = strings.TrimSpace(unitName)

		var response string
		if len(unitName) == 0 {
			response = m.Author.Mention() + ", the `!unit` command requires a unit card name. Please try again using this format `!unit <unit card name>`"
		} else {
			response = "```" + strings.Join(lookupUtil.LookupUnit(unitName), "\n") + "```"
		}

		s.ChannelMessageSend(m.ChannelID, response)
	}

	if strings.HasPrefix(m.Content, "!upgrade") {
		upgradeName := strings.Replace(m.Content, "!upgrade", "", 1)
		upgradeName = strings.TrimSpace(upgradeName)

		var response string
		if len(upgradeName) == 0 {
			response = m.Author.Mention() + ", the `!upgrade` command requires an upgrade card name. Please try again using this format `!upgrade <upgrade card name>`"
		} else {
			response = "```" + strings.Join(lookupUtil.LookupUpgrade(upgradeName), "\n") + "```"
		}

		s.ChannelMessageSend(m.ChannelID, response)
	}

	if strings.HasPrefix(m.Content, "!command") {
		commandName := strings.Replace(m.Content, "!command", "", 1)
		commandName = strings.TrimSpace(commandName)

		var response string
		if len(commandName) == 0 {
			response = m.Author.Mention() + ", the `!command` command requires a command card name. Please try again using this format `!command <command card name>`"
		} else {
			response = "```" + strings.Join(lookupUtil.LookupCommand(commandName), "\n") + "```"
		}

		s.ChannelMessageSend(m.ChannelID, response)
	}

	if strings.HasPrefix(m.Content, "!search") {
		searchText := strings.Replace(m.Content, "!search", "", 1)
		searchText = strings.TrimSpace(searchText)

		if searchText == "" {
			response := m.Author.Mention() + ", the `!search` command requires a search term. Please try again using this format `!search <search term>`"
			s.ChannelMessageSendEmbed(m.ChannelID, output.Error("Bad input", response))
		} else {
			embeddedResults := searchUtil.FullSearch(searchText)
			for _, embed := range embeddedResults {
				s.ChannelMessageSendEmbed(m.ChannelID, embed)
			}
		}
	}

}
