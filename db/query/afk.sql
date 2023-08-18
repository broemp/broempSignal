-- name: CreateAFK :one
INSERT INTO "afk" (
   userid
) VALUES (
    $1
) RETURNING *;

-- name: ListAFK :one
SELECT * from "afk"
WHERE userid = $1;

-- name: GetAFKCount :one
SELECT count(*) from "afk"
WHERE userid = $1;
