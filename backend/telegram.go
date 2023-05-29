package backend

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	Debug            = false
	bot              *tgbotapi.BotAPI
	codes            = make(map[int]int64)
	answeredMessages = make(map[int64]bool)
)

func InitTelegram() {
	var err error
	bot, err = tgbotapi.NewBotAPI(TelegramToken)
	rand.Seed(time.Now().UnixNano())

	if err != nil {
		log.Fatal("Couldn't Create Telegram Bot ", err)
	}

	bot.Debug = Debug
	TelegramBotName = bot.Self.UserName
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		var msg tgbotapi.MessageConfig
		if update.Message != nil && update.Message.IsCommand() {
			// ignore any non-command Messages
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				if !checkTelegramLink(update.Message.Chat.ID) {
					msg.Text = "Welcome to BroempSignal\nYour Code is: " + generateCode(update.Message.Chat.ID)
				} else {
					msg.Text = "You are already registered! If you want to register another account, please use /unregister first."
				}
			case "unregister":
				if !checkTelegramLink(update.Message.Chat.ID) {
					msg.Text = "You are not registered!"
				} else {
					msg.Text = "You are now unregistered!"
					removeTelegramLink(update.Message.Chat.ID)
				}
			default:
				continue
			}

			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}

		} else if update.CallbackQuery != nil {
			// Skip already processed messages
			if answeredMessages[int64(update.CallbackQuery.Message.MessageID)] {
				continue
			}

			// Answer callback query
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			switch update.CallbackQuery.Data {
			case "yes":
				acceptInvite(update.CallbackQuery.Message.Chat.ID)
				deleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)

				answeredMessages[int64(update.CallbackQuery.Message.MessageID)] = true
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "You are in!")
				var keyboard = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Cancel", "cancel"),
					))
				msg.ReplyMarkup = keyboard

				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}

			case "no":
				declineInvite(update.CallbackQuery.Message.Chat.ID)
				deleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				answeredMessages[int64(update.CallbackQuery.Message.MessageID)] = true

				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "You are out!\n Fuck you")
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}

			case "cancel":
				deleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Nice gehöppert!")
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			default:
				continue
			}
		}
	}
}

func generateCode(chatID int64) string {
	code := rand.Intn(999999-100000) + 100000
	codes[code] = chatID

	return strconv.Itoa(code)
}

func sendMessage(chatID int64, sender string) {
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Yes", "yes"),
			tgbotapi.NewInlineKeyboardButtonData("No", "no"),
		))
	msg := tgbotapi.NewMessage(chatID, "Assembling group of "+sender+" and his friends:\nAre You in?")
	msg.ReplyMarkup = keyboard
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func deleteMessage(chatId int64, messageId int) {
	// Delete old Message
	msgToDelete := tgbotapi.DeleteMessageConfig{
		ChatID:    chatId,
		MessageID: messageId,
	}
	_, err := bot.Request(msgToDelete)
	if err != nil {
		log.Println(err)
	}
}
