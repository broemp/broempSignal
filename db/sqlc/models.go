// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"database/sql"
)

type Group struct {
	Groupid   int64        `json:"groupid"`
	Hostid    int64        `json:"hostid"`
	Active    sql.NullBool `json:"active"`
	CreatedAt sql.NullTime `json:"created_at"`
}

type HoepperCount struct {
	Hoepper      int64         `json:"hoepper"`
	Victim       int64         `json:"victim"`
	HoepperCount sql.NullInt32 `json:"hoepperCount"`
}

type Invite struct {
	Groupid int64 `json:"groupid"`
	Guestid int64 `json:"guestid"`
	// null = waiting, false=declined, true=accepted
	Accepted  sql.NullBool `json:"accepted"`
	CreatedAt sql.NullTime `json:"created_at"`
}

type User struct {
	Userid     int64         `json:"userid"`
	Username   string        `json:"username"`
	Discordid  int64         `json:"discordid"`
	Telegramid sql.NullInt64 `json:"telegramid"`
	CreatedAt  sql.NullTime  `json:"created_at"`
}