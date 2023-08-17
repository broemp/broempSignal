package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var s *discordgo.Session

func InitDiscord(discordToken string, discordGuildId string) {
	log.Println("Starting Discord Bot")
	var err error
	s, err = discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatal("Couldn't create Discord session: ", err)
	}
	s.Open()
	registerCommands(s, discordGuildId)
}

// TODO: Refactor commands

var (
	integerOptionMinValue = 100000.0
	commands              = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Replies with pong",
		},
		{
			Name:        "afk",
			Description: "Report someone as afk",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "User",
					Description: "User who is afk",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionUser,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong!",
				},
			})
		},
	}
)

func registerCommands(s *discordgo.Session, discordGuildId string) {
	log.Println("Registering commands")

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, discordGuildId, commands)
	if err != nil {
		log.Fatal("Couldn't register commands: ", err)
	}
}

func CloseDiscord() {
	log.Println("Closing Discord")
	s.Close()
}

// TODO make embed or retain old accepts
func updateInteraction(interaction *discordgo.Interaction, msg string) {
	s.InteractionResponseEdit(interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}
