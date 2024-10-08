// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import ()

type BattleLog struct {
	ID       int64
	Position string
	Fecha    int64
}

type Grupo struct {
	ID     int64
	ChatID int64
}

type GrupoBattle struct {
	ID          int64
	GrupoID     int64
	BattleLogID int64
}

type Player struct {
	ID   int64
	Name string
	Team string
}

type PlayerBattle struct {
	ID          int64
	PlayerID    int64
	BattleLogID int64
}
