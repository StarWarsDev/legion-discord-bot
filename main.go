package main

import (
	"flag"
	"fmt"
	"github.com/StarWarsDev/legion-discord-bot/channel"
	"github.com/StarWarsDev/legion-discord-bot/search"
	"github.com/bwmarrin/discordgo"
	"github.com/shurcooL/graphql"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Hello, World! I am the Discord Legion bot!")

	// init all the data
	var token string
	var archivesURL string
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.StringVar(&archivesURL, "url", "https://sw-legion-archives.herokuapp.com/graphql", "Archives GraphQL URL")
	flag.Parse()

	if token == "" {
		// no token was assigned via the t flag, try looking in the environment
		token = os.Getenv("DISCORD_TOKEN")
	}

	if token == "" {
		// still no token, panic!
		panic("No discord token provided! Try passing it with the '-t' flag or setting 'DISCORD_TOKEN' in the environment.")
	}

	// create a new connection to Discord
	discord, err := discordgo.New("Bot " + token)

	if err != nil {
		panic(err)
	}

	// create the graphql client
	client := graphql.NewClient(archivesURL, nil)

	// create and add the message handler
	discord.AddHandler(channel.NewMessageHandler(client))

	// open the connection to Discord
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
	_ = discord.Close()
	_ = os.RemoveAll(search.IndexKey)
}
