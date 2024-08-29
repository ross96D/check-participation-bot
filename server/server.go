package server

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ross96D/battle-log-parser/parser"
	battleadmin "github.com/ross96D/cw_participation_bot/battle_admin"
	"github.com/ross96D/cw_participation_bot/tg"
	"github.com/rs/zerolog/log"
)

func Run(telegramTokenValue, serviceUrl string, poolling bool, port int) {
	if poolling {
		polling{telegramApiToken: telegramTokenValue, serviceUrl: serviceUrl}.runPolling()
	} else {
		runServer(telegramTokenValue, serviceUrl, port)
	}
}

func runServer(telegramTokenValue, serviceUrl string, port int) {
	s := Server(telegramTokenValue, serviceUrl)
	log.Debug().Msg("initializing server")

	cert, err := tls.LoadX509KeyPair("ca.crt", "ca.key")
	if err != nil {
		log.Panic().Err(err).Msg("loading tls cert")
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	server := http.Server{
		Addr:      ":" + strconv.FormatInt(int64(port), 10),
		Handler:   s,
		TLSConfig: config,
	}
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Error().Err(err).Send()
	}
}

func GetCommand(text string) (cmd string, ok bool) {
	if text == "" {
		return
	}
	if text[0] != '/' {
		return
	}
	i := 0
	for ; i < len(text); i++ {
		if text[i] == ' ' {
			break
		}
	}
	cmd, ok = text[0:i], true
	return
}

func Server(serviceUrl, telegramTokenValue string) *echo.Echo {
	s := echo.New()

	s.Use(NoError, Logger())

	s.POST("/hook", func(c echo.Context) error {
		var update tg.Update
		dec := json.NewDecoder(c.Request().Body)
		err := dec.Decode(&update)
		if err != nil {
			return err
		}
		cmd, ok := GetCommand(update.Message.Text)
		if !ok {
			return nil
		}

		switch cmd {
		case "/resumeByPlayer":
			webViewUrl, _ := strings.CutPrefix(update.Message.Text, cmd)
			webViewUrl = strings.TrimSpace(webViewUrl)
			pr, err := Parse(webViewUrl, serviceUrl)
			if err != nil {
				return err
			}
			SendPlayerResumeTg(telegramTokenValue, update.Message.Chat.ID, pr)

		case "upload", "/u":
			webViewUrl, _ := strings.CutPrefix(update.Message.Text, cmd)
			webViewUrl = strings.TrimSpace(webViewUrl)
			battle, err := GetBattle(webViewUrl, serviceUrl)
			if err != nil {
				return err
			}

			err = battleadmin.AddBattleToGroup(context.Background(), battle, update.Message.Chat.ID)
			if err != nil {
				return err
			}
			err = tg.New(telegramTokenValue).SendMessage(tg.SendMessage{
				Text:   "succesfully added " + webViewUrl,
				ChatID: update.Message.Chat.ID,
			})
			if err != nil {
				return err
			}

			return nil

		case "/playerParticipation", "/pp":
			playerName, _ := strings.CutPrefix(update.Message.Text, cmd)
			playerName = strings.TrimSpace(playerName)

			playerID, err := battleadmin.GetPlayerIDByName(context.Background(), playerName)
			if err != nil {
				if err == sql.ErrNoRows {
					return tg.New(telegramTokenValue).SendMessage(tg.SendMessage{
						Text:   "player not found",
						ChatID: update.Message.Chat.ID,
					})
				}
				return err
			}
			groupID, err := battleadmin.GroupByChatID(context.Background(), update.Message.Chat.ID)
			if err != nil {
				return err
			}

			countGroup, err := battleadmin.CountBattlesFromGroup(context.Background(), groupID)
			if err != nil {
				return err
			}
			countPlayer, err := battleadmin.CountBattlesFromPlayerAndGroup(context.Background(), groupID, playerID)
			if err != nil {
				return err
			}

			err = tg.New(telegramTokenValue).SendMessage(tg.SendMessage{
				Text:   fmt.Sprintf("group battles registered %d. Player battles %d", countGroup, countPlayer),
				ChatID: update.Message.Chat.ID,
			})
			if err != nil {
				return err
			}

			return nil

		default:
			return nil
		}
		return nil
	})

	return s
}

func Parse(webviewUrl string, serviceUrl string) ([]PlayerResume, error) {
	battle, err := GetBattle(webviewUrl, serviceUrl)
	if err != nil {
		return nil, err
	}

	resume := PlayerResumen(battle)

	b := strings.Builder{}

	ss := make([]PlayerResume, 0)
	for _, k := range sort(resume) {
		p := resume[k]
		b.WriteString(p.String())
		b.WriteByte('\n')
		ss = append(ss, p)
	}

	return ss, nil
}

func SendPlayerResumeTg(telegramTokenValue string, chatID int64, playersResume []PlayerResume) {
	lines := strings.Split(AllPlayerResume(playersResume).String(), "\n")

	for i := 0; i < len(lines); i++ {
		end := i + 80
		if len(lines)-1 < end {
			end = len(lines) - 1
		}
		toSend := strings.Join(lines[i:end], "\n")

		i = end

		err := tg.New(telegramTokenValue).SendMessage(tg.SendMessage{
			Text:      "```\n" + toSend + "```",
			ChatID:    chatID,
			ParseMode: "MarkdownV2",
		})
		if err != nil {
			log.Error().Err(err).Msgf("Send player resume %d", chatID)
			break
		}
	}

}

func GetBattle(webviewUrl string, serviceUrl string) (parser.Battle, error) {
	uri, err := url.Parse(serviceUrl)
	if err != nil {
		return parser.Battle{}, fmt.Errorf("url.Parse(%s) %w", serviceUrl, err)
	}
	queryValues := uri.Query()
	queryValues.Set("url", webviewUrl)
	uri.RawQuery = queryValues.Encode()

	resp, err := http.Get(uri.String())
	if err != nil {
		return parser.Battle{}, fmt.Errorf("Parse http.Get(serviceUrl) %w", err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return parser.Battle{}, fmt.Errorf("Parse io.ReadAll %w", err)
	}

	if resp.StatusCode != 200 {
		return parser.Battle{}, fmt.Errorf("http.Get(%s) status code %d %s", serviceUrl, resp.StatusCode, string(data))
	}

	battle := parser.Battle{}
	err = json.Unmarshal(data, &battle)
	if err != nil {
		return parser.Battle{}, fmt.Errorf("Parse json.Unmarshal(parser.Battle) %w", err)
	}
	return battle, nil
}
