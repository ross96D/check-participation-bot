-- name: GetAllPlayer :many
SELECT * FROM player;

-- name: GetIDByName :one
SELECT id FROM player WHERE name=?;

-- name: InsertIfNotExists :exec
INSERT OR REPLACE INTO player (name, team) 
    VALUES (?, ?);

-- name: CountBattlesFromPlayerAndGroup :one
SELECT count(pb.battle_log_id) FROM grupo_battle AS gb 
    JOIN player_battle AS pb ON pb.battle_log_id = gb.battle_log_id
    WHERE gb.grupo_id = ? AND pb.player_id = ?;

-- name: CountBattlesFromGroup :one
SELECT count(*) FROM grupo_battle AS gb 
    WHERE gb.grupo_id = ?;

