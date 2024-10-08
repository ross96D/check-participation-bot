// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: player.sql

package db

import (
	"context"
)

const countBattlesFromGroup = `-- name: CountBattlesFromGroup :one
SELECT count(*) FROM grupo_battle AS gb 
    WHERE gb.grupo_id = ?
`

func (q *Queries) CountBattlesFromGroup(ctx context.Context, grupoID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, countBattlesFromGroup, grupoID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countBattlesFromPlayerAndGroup = `-- name: CountBattlesFromPlayerAndGroup :one
SELECT count(pb.battle_log_id) FROM grupo_battle AS gb 
    JOIN player_battle AS pb ON pb.battle_log_id = gb.battle_log_id
    WHERE gb.grupo_id = ? AND pb.player_id = ?
`

type CountBattlesFromPlayerAndGroupParams struct {
	GrupoID  int64
	PlayerID int64
}

func (q *Queries) CountBattlesFromPlayerAndGroup(ctx context.Context, arg CountBattlesFromPlayerAndGroupParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, countBattlesFromPlayerAndGroup, arg.GrupoID, arg.PlayerID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getAllPlayer = `-- name: GetAllPlayer :many
SELECT id, name, team FROM player
`

func (q *Queries) GetAllPlayer(ctx context.Context) ([]Player, error) {
	rows, err := q.db.QueryContext(ctx, getAllPlayer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Player
	for rows.Next() {
		var i Player
		if err := rows.Scan(&i.ID, &i.Name, &i.Team); err != nil {
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

const getIDByName = `-- name: GetIDByName :one
SELECT id FROM player WHERE name=?
`

func (q *Queries) GetIDByName(ctx context.Context, name string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getIDByName, name)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const insertPlayerIfNotExists = `-- name: InsertPlayerIfNotExists :exec
INSERT OR IGNORE INTO player (name, team) 
    VALUES (?, ?)
`

type InsertPlayerIfNotExistsParams struct {
	Name string
	Team string
}

func (q *Queries) InsertPlayerIfNotExists(ctx context.Context, arg InsertPlayerIfNotExistsParams) error {
	_, err := q.db.ExecContext(ctx, insertPlayerIfNotExists, arg.Name, arg.Team)
	return err
}
