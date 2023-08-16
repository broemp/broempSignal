// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: invite.sql

package db

import (
	"context"
	"database/sql"
)

const createInvite = `-- name: CreateInvite :one
INSERT INTO "invite" (
   groupid, guestid
) VALUES (
    $1, $2
) RETURNING groupid, guestid, accepted, created_at
`

type CreateInviteParams struct {
	Groupid int64 `json:"groupid"`
	Guestid int64 `json:"guestid"`
}

func (q *Queries) CreateInvite(ctx context.Context, arg CreateInviteParams) (Invite, error) {
	row := q.db.QueryRowContext(ctx, createInvite, arg.Groupid, arg.Guestid)
	var i Invite
	err := row.Scan(
		&i.Groupid,
		&i.Guestid,
		&i.Accepted,
		&i.CreatedAt,
	)
	return i, err
}

const getInviteByGroupIdAndGuestId = `-- name: GetInviteByGroupIdAndGuestId :one
SELECT groupid, guestid, accepted, created_at FROM "invite"
WHERE groupid = $1 AND guestid = $2
`

type GetInviteByGroupIdAndGuestIdParams struct {
	Groupid int64 `json:"groupid"`
	Guestid int64 `json:"guestid"`
}

func (q *Queries) GetInviteByGroupIdAndGuestId(ctx context.Context, arg GetInviteByGroupIdAndGuestIdParams) (Invite, error) {
	row := q.db.QueryRowContext(ctx, getInviteByGroupIdAndGuestId, arg.Groupid, arg.Guestid)
	var i Invite
	err := row.Scan(
		&i.Groupid,
		&i.Guestid,
		&i.Accepted,
		&i.CreatedAt,
	)
	return i, err
}

const getInviteStatusByGroupIdAndGuestId = `-- name: GetInviteStatusByGroupIdAndGuestId :one
SELECT accepted FROM "invite"
WHERE groupid = $1 AND guestid = $2
`

type GetInviteStatusByGroupIdAndGuestIdParams struct {
	Groupid int64 `json:"groupid"`
	Guestid int64 `json:"guestid"`
}

func (q *Queries) GetInviteStatusByGroupIdAndGuestId(ctx context.Context, arg GetInviteStatusByGroupIdAndGuestIdParams) (sql.NullBool, error) {
	row := q.db.QueryRowContext(ctx, getInviteStatusByGroupIdAndGuestId, arg.Groupid, arg.Guestid)
	var accepted sql.NullBool
	err := row.Scan(&accepted)
	return accepted, err
}

const listInviteByGroupId = `-- name: ListInviteByGroupId :many
SELECT groupid, guestid, accepted, created_at FROM "invite"
WHERE groupid = $1
`

func (q *Queries) ListInviteByGroupId(ctx context.Context, groupid int64) ([]Invite, error) {
	rows, err := q.db.QueryContext(ctx, listInviteByGroupId, groupid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Invite{}
	for rows.Next() {
		var i Invite
		if err := rows.Scan(
			&i.Groupid,
			&i.Guestid,
			&i.Accepted,
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

const listInviteByGuest = `-- name: ListInviteByGuest :many
SELECT groupid, guestid, accepted, created_at from "invite"
WHERE guestid = $1
`

func (q *Queries) ListInviteByGuest(ctx context.Context, guestid int64) ([]Invite, error) {
	rows, err := q.db.QueryContext(ctx, listInviteByGuest, guestid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Invite{}
	for rows.Next() {
		var i Invite
		if err := rows.Scan(
			&i.Groupid,
			&i.Guestid,
			&i.Accepted,
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

const setInviteStatus = `-- name: SetInviteStatus :exec
UPDATE invite SET accepted=$3
WHERE groupid = $1 AND guestid = $2
`

type SetInviteStatusParams struct {
	Groupid  int64        `json:"groupid"`
	Guestid  int64        `json:"guestid"`
	Accepted sql.NullBool `json:"accepted"`
}

func (q *Queries) SetInviteStatus(ctx context.Context, arg SetInviteStatusParams) error {
	_, err := q.db.ExecContext(ctx, setInviteStatus, arg.Groupid, arg.Guestid, arg.Accepted)
	return err
}