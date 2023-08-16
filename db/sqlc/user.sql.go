// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO "user" (
    username, discordid, telegramid
) VALUES (
    $1, $2, $3
) RETURNING userid, username, discordid, telegramid, created_at
`

type CreateUserParams struct {
	Username   string        `json:"username"`
	Discordid  int64         `json:"discordid"`
	Telegramid sql.NullInt64 `json:"telegramid"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Discordid, arg.Telegramid)
	var i User
	err := row.Scan(
		&i.Userid,
		&i.Username,
		&i.Discordid,
		&i.Telegramid,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM "user"
WHERE userid = $1
`

func (q *Queries) DeleteUser(ctx context.Context, userid int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, userid)
	return err
}

const getUser = `-- name: GetUser :one
SELECT userid, username, discordid, telegramid, created_at from "user"
WHERE userid = $1
`

func (q *Queries) GetUser(ctx context.Context, userid int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, userid)
	var i User
	err := row.Scan(
		&i.Userid,
		&i.Username,
		&i.Discordid,
		&i.Telegramid,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByDiscordId = `-- name: GetUserByDiscordId :one
SELECT userid, username, discordid, telegramid, created_at FROM "user"
WHERE discordid = $1
`

func (q *Queries) GetUserByDiscordId(ctx context.Context, discordid int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByDiscordId, discordid)
	var i User
	err := row.Scan(
		&i.Userid,
		&i.Username,
		&i.Discordid,
		&i.Telegramid,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByTelegramId = `-- name: GetUserByTelegramId :one
SELECT userid, username, discordid, telegramid, created_at FROM "user"
WHERE telegramid = $1
`

func (q *Queries) GetUserByTelegramId(ctx context.Context, telegramid sql.NullInt64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByTelegramId, telegramid)
	var i User
	err := row.Scan(
		&i.Userid,
		&i.Username,
		&i.Discordid,
		&i.Telegramid,
		&i.CreatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT userid, username, discordid, telegramid, created_at from "user"
ORDER BY userid
LIMIT $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.Userid,
			&i.Username,
			&i.Discordid,
			&i.Telegramid,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTelegramId = `-- name: UpdateTelegramId :exec
UPDATE "user" SET telegramid = $2
WHERE userid = $1
`

type UpdateTelegramIdParams struct {
	Userid     int64         `json:"userid"`
	Telegramid sql.NullInt64 `json:"telegramid"`
}

func (q *Queries) UpdateTelegramId(ctx context.Context, arg UpdateTelegramIdParams) error {
	_, err := q.db.ExecContext(ctx, updateTelegramId, arg.Userid, arg.Telegramid)
	return err
}