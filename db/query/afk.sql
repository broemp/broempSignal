-- name: CreateAFK :one
INSERT INTO "afk" (
   userid
) VALUES (
    $1
) RETURNING *;

-- name: ListAFK :many
SELECT * from "afk"
WHERE userid = $1;

-- name: GetAFKCount :one
SELECT count(*) from "afk"
WHERE userid = $1;

-- name: ListUsersByAFKCount :many
SELECT userid, username, count(*) from "afk"
JOIN "user" ON "user".discordid = "afk".userid 
GROUP BY userid, username
ORDER BY count(*) desc
LIMIT 15;
