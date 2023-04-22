package backend

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	Debug = false
	bot   *tgbotapi.BotAPI
	codes = make(map[int]int64)
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

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() { // ignore any non-command Messages
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				msg.Text = "Welcome to Broempsignal! Please use /register to register your account."
			case "register":
				msg.Text = "Your Code is: " + generateCode(update.Message.Chat.ID)
			}
		}

		if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "start":
			msg.Text = "Welcome to Broempsignal! Please use /register to register your account."
		case "register":
			msg.Text = "Your Code is: " + generateCode(update.Message.Chat.ID)
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func generateCode(chatID int64) string {
	code := rand.Intn(999999-100000) + 100000
	codes[code] = chatID

	return strconv.Itoa(code)
}

var response = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Yes", "yes"),
		tgbotapi.NewInlineKeyboardButtonData("No", "no"),
	))

func sendMessage(chatID int64, sender string) {
	msg := tgbotapi.NewMessage(chatID, "Assembling group of "+sender+" and his friends:\nAre You in?")
	msg.ReplyMarkup = response
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}
