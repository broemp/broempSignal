-- name: StartHoepperCount :one
INSERT INTO "hoepperCount" (
    hoepper, victim, count
) VALUES (
    $1, $2, 1
) RETURNING *;
