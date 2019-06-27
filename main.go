package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token string
)

func init() {

	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
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
		s.ChannelMessageSend(m.ChannelID, "Sorry, "+m.Author.Mention()+"I can't help you right now since my creator didn't give me any other commands.\nI lied, you can always `!gonk` for fun.")
	}

	if m.Content == "!gonk" {
		s.ChannelMessageSend(m.ChannelID, "GONK")
	}
}
