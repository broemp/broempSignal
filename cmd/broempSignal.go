package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/broemp/broempSignal/api"
	"github.com/broemp/broempSignal/bots/discord"
	db "github.com/broemp/broempSignal/db/sqlc"
	"github.com/broemp/broempSignal/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}

	if config.DiscordToken != "" {
		discord.InitDiscord(config.DiscordToken, config.DiscordGuildId)
	}
}
