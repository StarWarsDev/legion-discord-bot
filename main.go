package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token      string
	legionData *LegionData
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()

	legionData = loadLegionData()
}

func main() {
	fmt.Println("Hello, World! I am the Discord Legion bot!")

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
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!help" {
		helpText := []string{
			"!unit <unit card name> - Displays information about the specified unit (WIP)",
			"!upgrade <upgrade card name> - Displays information about the specified upgrade (WIP)",
			"!command <command card name> - Displays information about the specified command card (WIP)",
			"!gonk - GONK",
			"!help - This help message",
		}
		s.ChannelMessageSend(m.ChannelID, "```"+strings.Join(helpText, "\n")+"```")
	}

	if m.Content == "!gonk" {
		s.ChannelMessageSend(m.ChannelID, "GONK")
	}

	if strings.HasPrefix(m.Content, "!unit") {
		unitName := strings.Replace(m.Content, "!unit", "", 1)
		unitName = strings.TrimSpace(unitName)

		var response string
		if len(unitName) == 0 {
			response = m.Author.Mention() + ", the `!unit` command requires a unit card name. Please try again using this format `!unit <unit card name>`"
		} else {
			response = "```" + strings.Join(lookupUnit(unitName), "\n") + "```"
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
			response = "```" + strings.Join(lookupUpgrade(upgradeName), "\n") + "```"
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
			response = "```" + strings.Join(lookupCommand(commandName), "\n") + "```"
		}

		s.ChannelMessageSend(m.ChannelID, response)
	}
}

func lookupUnit(unitName string) []string {
	cardInfo := []string{"Your mom is " + unitName}

	return cardInfo
}

func lookupUpgrade(upgradeName string) []string {
	cardInfo := []string{"Your dad has " + upgradeName}

	return cardInfo
}

func lookupCommand(commandName string) []string {
	cardInfo := []string{}

	for _, card := range legionData.CommandCards {
		if strings.ToLower(card.Name) == strings.ToLower(commandName) {

			/*
			   {
			   	"ldf":"zxflameprojector",
			   	"name":"ZX Flame Projector",
			   	"pips":2,
			   	"orders":"Boba Fett",
			   	"description":"During Boba Fett's activation he gains the following weapons:",
			   	"weapon":{
			   		"range":{
			   			"from":0,
			   			"to":1
			   		},
			   		"dice":{
			   			"red":1
			   		},
			   		"keywords":["Blast","Spray"]
			   	}
			   }
			*/

			cardInfo = []string{
				"Name: " + card.Name,
				fmt.Sprintf("Pips: %d", card.Pips),
				"Orders: " + card.Orders,
			}

			if len(card.Description) > 0 {
				cardInfo = append(cardInfo, "Description: "+card.Description)
			}

			diceCount := card.Weapon.Dice.White + card.Weapon.Dice.Black + card.Weapon.Dice.Red
			if diceCount > 0 {
				weaponInfo := []string{"Weapon:"}
				if len(card.Weapon.Name) > 0 {
					weaponInfo = append(weaponInfo, "  "+card.Weapon.Name)
				}

				weaponInfo = append(weaponInfo, fmt.Sprintf("  Range: %d - %d", card.Weapon.Range.From, card.Weapon.Range.To))
				weaponInfo = append(weaponInfo, fmt.Sprintf("  Dice: white: %d, black: %d, red: %d", card.Weapon.Dice.White, card.Weapon.Dice.Black, card.Weapon.Dice.Red))

				if len(card.Weapon.Keywords) > 0 {
					keywords := strings.Join(card.Weapon.Keywords, ", ")
					weaponInfo = append(weaponInfo, "  Keywords: "+keywords)
				}

				cardInfo = append(cardInfo, weaponInfo...)
			}

			return cardInfo
		}
	}

	if len(cardInfo) == 0 {
		cardInfo = []string{"Nothing found for \"" + commandName + "\""}
	}

	return cardInfo
}
