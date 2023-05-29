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
			user.Friends[i] = user.Friends[len(user.Friends)-1] // Copy last element to index i.
			user.Friends[len(user.Friends)-1] = model.User{}    // Erase last element (write zero value).
			user.Friends = user.Friends[:len(user.Friends)-1]   // Truncate slice.
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

func registerUser(user model.User, telegramCode int) error {
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

// Returns True if telegram is connected to discord
func checkTelegramLink(userId int64) bool {
	var user model.User
	result := db.First(&user, "telegram_id = ?", userId)
	return result.Error == nil
}

//Remove Telegram Link
func removeTelegramLink(userId int64) bool {
	var user model.User
	result := db.First(&user, "telegram_id = ?", userId)

	if result.Error != nil {
		return false
	}

	user.TelegramId = nil
	db.Save(&user)
	return true
}

func removeAllUserData(discordID string) bool {
	var user model.User
	result := db.Unscoped().Where("discord_id = ?", discordID).Delete(&user)
	if Debug {
		log.Println(result)
	}
	return result.Error == nil && result.RowsAffected > 0
}
