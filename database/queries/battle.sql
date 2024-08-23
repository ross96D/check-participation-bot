-- name: InsertBattle :one
INSERT INTO battle_log (position, fecha) VALUES (?, ?) RETURNING id;
