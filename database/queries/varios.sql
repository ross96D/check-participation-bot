-- name: InsertPlayerBattle :exec 
INSERT INTO player_battle (player_id, battle_log_id) VALUES (?, ?);

-- name: InsertGroupBattle :exec 
INSERT INTO grupo_battle (grupo_id, battle_log_id) VALUES (?, ?);

-- name: GroupByChatID :one
SELECT * FROM grupo WHERE chat_id = ?;

-- name: InsertGroup :exec
INSERT OR IGNORE INTO grupo (chat_id) 
    VALUES (?);


-- name: GetAllGroups :many
SELECT id FROM grupo;
