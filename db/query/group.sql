-- name: CreateGroup :one
INSERT INTO "group" (
    hostid
) VALUES (
    $1
) RETURNING *;

-- name: GetGroupByHost :one
SELECT * FROM "group"
WHERE hostid = $1;

-- name: ListGroupByGuest :many
SELECT * FROM "group"
NATURAL JOIN "invite"
WHERE guestid = $1;
