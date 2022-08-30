package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/StarWarsDev/legion-discord-bot/internal/channel"
	"github.com/StarWarsDev/legion-discord-bot/internal/data"
	"github.com/bwmarrin/discordgo"
)

func main() {
	log.Println("Hello, World! I am the Discord Legion bot!")

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
	s, err := discordgo.New("Bot " + token)

	if err != nil {
		panic(err)
	}

	inMemClient := data.InMemoryClient{}
	go func() {
		for {
			inMemClient.FetchAllData()
			time.Sleep(time.Minute * 5)
		}
	}()

	// create the graphql client
	client := data.NewArchivesClient(archivesURL)

	// create and add the message handler
	s.AddHandler(channel.NewMessageHandler(&client))

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "legion",
			Description: "Main command to start getting Legion data",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "command",
					Description: "Look up command info",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "term",
							Description: "Search term",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "field",
							Description: "Field to search",
							Required:    false,
						},
					},
				},
				{
					Name:        "keyword",
					Description: "Look up cards by keyword",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "term",
							Description: "Search term",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "field",
							Description: "Field to search",
							Required:    false,
						},
					},
				},
				{
					Name:        "unit",
					Description: "Look up unit info",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "term",
							Description: "Search term",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "field",
							Description: "Field to search",
							Required:    false,
						},
					},
				},
				{
					Name:        "upgrade",
					Description: "Look up upgrade info",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "term",
							Description: "Search term",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "field",
							Description: "Field to search",
							Required:    false,
						},
					},
				},
			},
		},
	}

	handlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"legion": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			var (
				// content string = ""
				command string = options[0].Name
				field   string = "name"
				term    string = options[0].Options[0].StringValue()
			)

			if len(options[0].Options) > 1 {
				field = options[0].Options[1].StringValue()
			}

			channel.HandleSlashCommand(command, field, term, s, i, &client, &inMemClient)
		},
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := handlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	// open the connection to Discord
	err = s.Open()
	if err != nil {
		panic(err)
	}

	// set the bot status
	// _ = discord.UpdateStatus(0, "Type !help to start")

	// add commands
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: '%v'", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	for _, cmd := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: '%v'", cmd.Name, err)
		}
	}

	// Cleanly close down the Discord session.
	_ = s.Close()
}
