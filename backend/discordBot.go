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
			Name:        "data",
			Description: "Shows what data Broemp Signal stores about you",
		},
		{
			Name:        "unregister",
			Description: "Unregister your account/Remove all your data",
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
		"data": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Broemp Signal stores the following data:\n\n" +
						"Your Discord ID\n" +
						"Your Discord Username\n" +
						"The Telegram Chat ID with the Bot(Not your phone number!)\n" +
						"Your Friends Discord ID\n\n" +
						"Broemp Signal does not store any messages or any other data\n" +
						"Broemp Signal does not share any data with any third party\n" +
						"You can delete your data at any time with /unregister\n\n",
				},
			})
		},
		"unregister": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if removeAllUserData(i.Member.User.ID) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Removed all of your Data!\nIf you want to use Broemp Signal again, you have to register again!",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You are not registered!",
					},
				})
			}
		},
		"assemble": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			msg := assemble(i.Member.User.ID, i.Interaction)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
		},
		"register": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if len(i.ApplicationCommandData().Options) == 0 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Please start a Conversation in the Telegram App with the Bot @" + TelegramBotName + "\n\n" +
							"You can get a code by sending /start to the Telegram Bot\n" +
							"Afterwards you can connect your account to Discord with /register <code>\n\n" +
							"You can start the chat by clicking this link on your phone https://t.me/" + TelegramBotName +
							"\n\n Or use Telegram Web at https://web.telegram.org/#/im?p=@" + TelegramBotName,
					},
				})
			} else {
				user := model.User{
					Name:      i.Member.User.Username,
					DiscordId: i.Member.User.ID,
				}

				code := int(math.Round(i.ApplicationCommandData().Options[0].Value.(float64)))

				err := registerUser(user, code)

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

// TODO make embed or retain old accepts
func updateInteraction(interaction *discordgo.Interaction, msg string) {
	s.InteractionResponseEdit(interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}
