-- name: CreateUser :one
INSERT INTO "user" (
    discordid, username, telegramid
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUser :one
SELECT * FROM "user"
WHERE discordid = $1;

-- name: GetUserByTelegramId :one
SELECT * FROM "user"
WHERE telegramid = $1;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE discordid = $1;

-- name: UpdateTelegramId :exec
UPDATE "user" SET telegramid = $2
WHERE discordid = $1;

-- name: ListUsers :many
SELECT * from "user"
ORDER BY discordid
LIMIT $1
OFFSET $2;
