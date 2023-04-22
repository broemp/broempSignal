package backend

import (
	"errors"
	"fmt"
	"log"

	"github.com/broemp/broempsignal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb() {
	log.Println("InitDb")
	var err error
	db, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	db.AutoMigrate(&model.User{})
}

func AddFriend(userId string, friend model.User) string {
	log.Println(userId + " Added friend " + friend.DiscordId)
	if friend.DiscordId == userId {
		return "Can't add yourself!"
	}
	var user model.User
	result := db.Preload("Friends").First(&user, "discord_id = ?", userId)
	if result.Error != nil {
		db.Save(&model.User{DiscordId: userId})
	}

	for _, f := range user.Friends {
		if f.DiscordId == friend.DiscordId {
			return "Already friends!"
		}
	}
	var tmp model.User
	result = db.First(&tmp, "discord_id = ?", friend.DiscordId)
	if result.Error != nil {
		db.Create(&friend)
	}

	user.Friends = append(user.Friends, friend)
	db.Save(&user)

	return "Added " + friend.Name + " as friend!"
}

func RemoveFriend(userId string, friendId string) string {
	var user model.User
	db.Preload("Friends").First(&user, "discord_id = ?", userId)

	for i, f := range user.Friends {
		if f.DiscordId == friendId {
			user.Friends = append(user.Friends[:i], user.Friends[i+1:]...)
			db.Save(&user)
			return "Removed " + f.Name + " as friend!"
		}
	}
	return "Couldn't find friend!"
}

func ListFriends(userId string) string {
	var user model.User
	db.Preload("Friends").First(&user, "discord_id = ?", userId)
	var msg string = "Your friends:\n"
	for _, f := range user.Friends {
		msg += f.Name + "\n"
	}
	return msg
}

func RegisterUser(user model.User, telegramCode int) error {
	log.Println("Registered user " + user.Name + " with code " + fmt.Sprint(telegramCode))

	chatId, ok := codes[telegramCode]
	if !ok {
		return errors.New("code not found")
	}
	delete(codes, telegramCode)

	result := db.First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user.TelegramId = &chatId

		// Create new User if no user with this discordId exists
		result = db.Create(&user)
		if result.Error != nil {
			log.Println(result.Error.Error())
			return errors.New("couldn't create user")
		}
		return nil
	} else if result.Error != nil {
		return result.Error
	}

	//Update User
	user.TelegramId = &chatId
	result = db.Model(&user).Updates(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func assemble(userId string) string {
	var user model.User
	db.Preload("Friends").First(&user, "discord_id = ?", userId)
	for _, f := range user.Friends {
		sendMessage(*f.TelegramId, user.Name)
	}
	return "You have no friends!"
}
