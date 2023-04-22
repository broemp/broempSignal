package backend

import (
	"fmt"
	"log"
	"math"

	"github.com/broemp/broempsignal/model"
	"github.com/bwmarrin/discordgo"
)

var (
	DiscordToken    string
	GuildId         string
	TelegramToken   string
	TelegramBotName string
)

var s *discordgo.Session

func InitDiscord() {
	log.Println("Starting Discord Bot")
	var err error
	s, err = discordgo.New("Bot " + DiscordToken)
	if err != nil {
		log.Fatal("Couldn't create Discord session: ", err)
	}
	s.Open()
	registerCommands(s)
}

var (
	integerOptionMinValue = 100000.0
	commands              = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Replies with pong",
		},
		{
			Name:        "assemble",
			Description: "Assembles a group of people",
		},
		{
			Name:        "friends",
			Description: "Manage your friends",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "add",
					Description: "Add a friend",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "user",
							Description: "The user to add",
							Required:    true,
							Type:        discordgo.ApplicationCommandOptionUser,
						},
					}},
				{
					Name:        "remove",
					Description: "Remove a friend",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "user",
							Description: "The user to Remove",
							Required:    true,
							Type:        discordgo.ApplicationCommandOptionUser,
						},
					}},
				{
					Name:        "list",
					Description: "Get a list of your friends",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:        "register",
			Description: "Register your account with your Telegram account",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "telegram_code",
					Description: "The Code you got from the Telegram Bot",
					Required:    false,
					Type:        discordgo.ApplicationCommandOptionInteger,
					MinValue:    &integerOptionMinValue,
					MaxValue:    999999,
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
		"assemble": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			assemble(i.Member.User.ID)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Assembling...",
				},
			})
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Found 3 people!",
				},
			})
		},
		"register": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if len(i.ApplicationCommandData().Options) == 0 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Please start a Conversation with the Telegram Bot @" + TelegramBotName +
							" and send the code you get from it here\nYou can get the code by sending /start in the Telegram Bot\nAfter that you can register your account on Discord with /register <code>",
					},
				})
			} else {
				user := model.User{
					Name:      i.Member.User.Username,
					DiscordId: i.Member.User.ID,
				}

				code := int(math.Round(i.ApplicationCommandData().Options[0].Value.(float64)))

				err := RegisterUser(user, code)

				if err == nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Registered your account",
						},
					})
				} else if err.Error() == "couldn't create user" {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "User already exists",
						},
					})
				} else if err.Error() == "code not found" {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Wrong Code. Please try again",
						},
					})
				}
			}
		},
		"friends": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			content := ""

			switch options[0].Name {
			case "add":
				friend := model.User{
					DiscordId: options[0].Options[0].UserValue(s).ID,
					Name:      options[0].Options[0].UserValue(s).Username,
				}
				fmt.Println(friend.DiscordId)
				fmt.Println(friend.Name)
				content = AddFriend(i.Member.User.ID, friend)
			case "remove":
				content = RemoveFriend(i.Member.User.ID, options[0].Options[0].UserValue(s).ID)
			case "list":
				content = ListFriends(i.Member.User.ID)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
	}
)

func registerCommands(s *discordgo.Session) {
	log.Println("Registering commands")

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, GuildId, commands)
	if err != nil {
		log.Fatal("Couldn't register commands: ", err)
	}

}

func CloseDiscord() {
	log.Println("Closing Discord")
	s.Close()
}
