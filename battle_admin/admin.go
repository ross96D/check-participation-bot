package battleadmin

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	"github.com/ross96D/battle-log-parser/parser"
	"github.com/ross96D/cw_participation_bot/database/connection"
	"github.com/ross96D/cw_participation_bot/database/db"
	"github.com/rs/zerolog/log"
)

func GetPlayers(battle parser.Battle) []parser.User {
	result := []parser.User{}
	// complejidad n^2?
	for _, turn := range battle.Turns {
		user := turn.Attacker
		if !user.IsMiss() && !slices.Contains(result, user) {
			result = append(result, user)
		} else {
			user = turn.Target
			if !user.IsMiss() && !slices.Contains(result, user) {
				result = append(result, user)
			}
		}
	}
	return result
}

func CheckGroup(ctx context.Context, chatID int64) (bool, error) {
	conn := connection.DB()
	_, err := db.New(conn).GroupByChatID(ctx, chatID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func AddBattleToGroup(ctx context.Context, battle parser.Battle, chatID int64) error {
	tx, err := connection.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Add Battle BeginTx %w", err)
	}
	if err = addBattleToGroup(ctx, tx, battle, chatID); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GroupByChatID(ctx context.Context, chatID int64) (int64, error) {
	return groupByChatID(ctx, connection.DB(), chatID)
}

func groupByChatID(ctx context.Context, conn db.DBTX, chatID int64) (int64, error) {
	group, err := db.New(conn).GroupByChatID(ctx, chatID)
	if err != nil {
		return 0, fmt.Errorf("Add BattleGroup GroupByChatID %w", err)
	}
	return group.ID, nil
}

func insertGroup(ctx context.Context, conn db.DBTX, chatID int64) (int64, error) {
	return db.New(conn).InsertGroup(ctx, chatID)
}

func addBattleToGroup(ctx context.Context, conn db.DBTX, battle parser.Battle, chatID int64) error {
	battleID, err := addBattle(ctx, conn, battle)
	if err != nil {
		return err
	}

	groupID, err := insertGroup(ctx, conn, chatID)
	if err != nil {
		return fmt.Errorf("Add BattleGroup insertGroup %w", err)
	}

	err = db.New(conn).InsertGroupBattle(ctx, db.InsertGroupBattleParams{
		GrupoID:     groupID,
		BattleLogID: battleID,
	})
	if err != nil {
		return fmt.Errorf("Add Battle InsertGroupBattle groupid %d battleid %d %w", chatID, battleID, err)
	}
	return nil
}

func AddBattle(ctx context.Context, battle parser.Battle) error {
	tx, err := connection.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Add Battle BeginTx %w", err)
	}

	_, err = addBattle(ctx, tx, battle)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func addBattle(ctx context.Context, conn db.DBTX, battle parser.Battle) (int64, error) {
	playersIds, err := insertPlayers(ctx, conn, battle)
	if err != nil {
		return 0, fmt.Errorf("Add Battle insertPlayers %w", err)
	}

	var battleID int64
	battleID, err = insertBattle(ctx, conn, battle)
	if err != nil {
		return 0, fmt.Errorf("Add Battle insertBattle %w", err)
	}

	err = insertPlayerBattle(ctx, conn, battleID, playersIds)
	if err != nil {
		return 0, fmt.Errorf("Add Battle insertPlayerBattle %w", err)
	}
	return battleID, nil
}

func insertPlayers(ctx context.Context, conn db.DBTX, battle parser.Battle) (ids []int64, err error) {
	ids = make([]int64, 0)
	for _, player := range GetPlayers(battle) {
		var id int64
		id, err = db.New(conn).InsertIfNotExists(ctx, db.InsertIfNotExistsParams{
			Name: player.Name,
			Team: player.Team.String(),
		})
		if err != nil {
			return
		}
		ids = append(ids, id)
	}
	return
}

func insertBattle(ctx context.Context, conn db.DBTX, battle parser.Battle) (id int64, err error) {
	id, err = db.New(conn).InsertBattle(ctx, db.InsertBattleParams{
		Position: battle.Resume.Position.String(),
		Fecha:    battle.Date.UnixMilli(),
	})
	if err != nil {
		return
	}
	row := conn.QueryRowContext(ctx, "select id, fecha from battle_log where id=?", id)
	if row.Err() != nil {
		log.Error().Err(row.Err()).Send()
	} else {
		dd := db.BattleLog{}
		err = row.Scan(&dd.ID, &dd.Fecha)
		if err != nil {
			log.Error().Err(err).Send()
		}
	}

	return
}

func insertPlayerBattle(ctx context.Context, conn db.DBTX, battleID int64, playersID []int64) (err error) {
	for _, playerID := range playersID {
		err = db.New(conn).InsertPlayerBattle(ctx, db.InsertPlayerBattleParams{
			BattleLogID: battleID,
			PlayerID:    playerID,
		})
		if err != nil {
			return err
		}
	}
	return
}

func CountBattlesFromGroup(ctx context.Context, groupID int64) (int64, error) {
	return db.New(connection.DB()).CountBattlesFromGroup(ctx, groupID)
}

func CountBattlesFromPlayerAndGroup(ctx context.Context, groupID int64, playerID int64) (int64, error) {
	return db.New(connection.DB()).CountBattlesFromPlayerAndGroup(ctx, db.CountBattlesFromPlayerAndGroupParams{
		GrupoID:  groupID,
		PlayerID: playerID,
	})
}

func GetPlayerIDByName(ctx context.Context, name string) (int64, error) {
	return db.New(connection.DB()).GetIDByName(ctx, name)
}
