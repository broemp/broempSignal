package model

type User struct {
	Name       string
	DiscordId  string `gorm:"primaryKey"`
	TelegramId *int64
	Friends    []User `gorm:"many2many:user_friends;"`
}
