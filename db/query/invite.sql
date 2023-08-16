-- name: CreateInvite :one
INSERT INTO "invite" (
   groupid, guestid
) VALUES (
    $1, $2
) RETURNING *;

-- name: ListInviteByGroupId :many
SELECT * FROM "invite"
WHERE groupid = $1;

-- name: ListInviteByGuest :many
SELECT * from "invite"
WHERE guestid = $1;

-- name: GetInviteByGroupIdAndGuestId :one
SELECT * FROM "invite"
WHERE groupid = $1 AND guestid = $2;

-- name: GetInviteStatusByGroupIdAndGuestId :one
SELECT accepted FROM "invite"
WHERE groupid = $1 AND guestid = $2;

-- name: SetInviteStatus :exec
UPDATE invite SET accepted=$3
WHERE groupid = $1 AND guestid = $2;
