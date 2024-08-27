package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	battleadmin "github.com/ross96D/cw_participation_bot/battle_admin"
	"github.com/ross96D/cw_participation_bot/tg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type polling struct {
	telegramApiToken string
	serviceUrl       string
}

func (p polling) runPolling() {
	var updates []tg.Update
	var err error
	var retryDelay time.Duration
	var offset int64
	for {
		updates, err = tg.GetUpdates(p.telegramApiToken, tg.GetUpdatesParams{
			Offset:         offset,
			Timeout:        uint(time.Second),
			AllowedUpdates: []tg.AllowedUpdate{tg.AUMessage},
		})
		if err != nil {
			log.Error().Err(err).Msg("GetUpdates")
			if retryDelay == 0 {
				retryDelay = time.Second
			} else {
				retryDelay *= 2
			}
			time.Sleep(retryDelay)
			continue
		} else {
			retryDelay = 0
		}
		max := int64(0)
		for _, u := range updates {
			if u.ID > max {
				max = u.ID
			}
			go p.handleUpdate(u)
		}
		offset = max + 1
		time.Sleep(time.Millisecond * 100)
	}
}

func (p polling) handleUpdate(u tg.Update) {
	start := time.Now()
	ok, err := p.handle(u)
	if !ok {
		return
	}
	var logger *zerolog.Event
	if err != nil {
		logger = log.Error()
		logger.Err(err)
	} else {
		logger = log.Info()
	}
	logger.Int64("update_id", u.ID).Str("message", u.Message.Text).Dur("elapsed", time.Since(start))
	logger.Send()
}

func (p polling) handle(update tg.Update) (bool, error) {
	cmd, ok := GetCommand(update.Message.Text)
	if !ok {
		return false, nil
	}
	switch cmd {
	case "/resumeByPlayer":
		webViewUrl, _ := strings.CutPrefix(update.Message.Text, cmd)
		webViewUrl = strings.TrimSpace(webViewUrl)
		pr, err := Parse(webViewUrl, p.serviceUrl)
		if err != nil {
			return true, err
		}
		SendPlayerResumeTg(p.telegramApiToken, update.Message.Chat.ID, pr)

		return true, nil

	case "upload", "/u":
		webViewUrl, _ := strings.CutPrefix(update.Message.Text, cmd)
		webViewUrl = strings.TrimSpace(webViewUrl)
		battle, err := GetBattle(webViewUrl, p.serviceUrl)
		if err != nil {
			return true, err
		}

		ok, err := battleadmin.CheckGroup(context.Background(), update.Message.Chat.ID)
		if err != nil {
			return true, err
		}
		if !ok {
			return true, errors.New("CheckGroup false " + strconv.FormatInt(update.Message.Chat.ID, 10))
		}

		err = battleadmin.AddBattleToGroup(context.Background(), battle, update.Message.Chat.ID)
		if err != nil {
			return true, err
		}
		err = tg.New(p.telegramApiToken).SendMessage(tg.SendMessage{
			Text:   "succesfully added " + webViewUrl,
			ChatID: update.Message.Chat.ID,
		})
		if err != nil {
			return true, err
		}

		return true, nil

	case "/playerParticipation", "/pp":
		playerName, _ := strings.CutPrefix(update.Message.Text, cmd)
		playerName = strings.TrimSpace(playerName)

		playerID, err := battleadmin.GetPlayerIDByName(context.Background(), playerName)
		if err != nil {
			if err == sql.ErrNoRows {
				return true, tg.New(p.telegramApiToken).SendMessage(tg.SendMessage{
					Text:   "player not found",
					ChatID: update.Message.Chat.ID,
				})
			}
			return true, err
		}
		groupID, err := battleadmin.GroupByChatID(context.Background(), update.Message.Chat.ID)
		if err != nil {
			return true, err
		}

		countGroup, err := battleadmin.CountBattlesFromGroup(context.Background(), groupID)
		if err != nil {
			return true, err
		}
		countPlayer, err := battleadmin.CountBattlesFromPlayerAndGroup(context.Background(), groupID, playerID)
		if err != nil {
			return true, err
		}

		err = tg.New(p.telegramApiToken).SendMessage(tg.SendMessage{
			Text:   fmt.Sprintf("group battles registered %d. Player battles %d", countGroup, countPlayer),
			ChatID: update.Message.Chat.ID,
		})
		if err != nil {
			return true, err
		}

		return true, nil

	default:
		return false, nil
	}
}
