package main

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/broemp/broempsignal/backend"
	"github.com/joho/godotenv"
)

func main() {
	// Load Environment Variables
	godotenv.Load()
	discordToken, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok || discordToken == "" {
		log.Fatal("DISCORD_TOKEN not set in env")
	}
	guildId, ok := os.LookupEnv("GUILD_ID")
	if !ok || guildId == "" {
		log.Fatal("GUILD_ID not set in env")
	}
	telegramToken, ok := os.LookupEnv("TELEGRAM_TOKEN")
	if !ok || telegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN not set in env")
	}
	debug, _ := os.LookupEnv("DEBUG")
	if strings.ToLower(debug) == "true" {
		backend.Debug = true
	}

	backend.DiscordToken = discordToken
	backend.GuildId = guildId
	backend.TelegramToken = telegramToken

	backend.InitDb()
	backend.InitDiscord()
	backend.InitTelegram()

	defer backend.CloseDiscord()

	// Wait here until CTRL-C or other term signal is received.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
