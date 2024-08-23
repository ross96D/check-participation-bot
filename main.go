package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"flag"

	"github.com/ross96D/cw_participation_bot/database/connection"
	"github.com/ross96D/cw_participation_bot/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var telegramTokenValue string
var tursoTokenValue string

var port uint
var tokenTg string
var serviceUrl string
var databaseName string
var tokenTurso string

var cliWebviewUrl string

func main() {
	flag.UintVar(&port, "port", 0, "port to listen on")
	flag.StringVar(&tokenTg, "token-tg", "", "telegram api token value")
	flag.StringVar(&tokenTg, "token-turso", "", "turso token value")
	flag.StringVar(&serviceUrl, "service", "", "service url for parsing web view")
	flag.StringVar(&cliWebviewUrl, "req", "", "battle log webview url")
	flag.StringVar(&databaseName, "db", "", "database url ex (dbname-org.turso.io)")

	flag.Parse()

	if port == 0 {
		println("required port value to be set (--port=<number>)")
		os.Exit(1)
	}

	if databaseName == "" {
		println("required database url to be set (--db=<url> ex: --db=dbname-org.turso.io)")
		os.Exit(1)
	}

	if serviceUrl == "" {
		println("required service url to be set (--service=<URL>)")
		os.Exit(1)
	}

	var ok bool

	if tokenTg != "" {
		telegramTokenValue = tokenTg
	} else {
		// default look up for CHK_PART_TELEGRAM_API_TOKEN
		telegramTokenValue, ok = os.LookupEnv("CW_BOT_TELEGRAM_API_TOKEN")
		if !ok {
			fmt.Println(os.Environ())
			println("no telegram api token value provided, default to enviroment variable CW_BOT_TELEGRAM_API_TOKEN	 but value was not found")
			os.Exit(1)
		}
	}

	if tokenTurso != "" {
		tursoTokenValue = tokenTurso
	} else {
		tursoTokenValue, ok = os.LookupEnv("CW_BOT_TURSO_TOKEN")
		if !ok {
			println("no turso token value provided, default to enviroment variable CW_BOT_TURSO_TOKEN but value was not found")
			os.Exit(1)
		}
	}

	log.Logger = log.Output(zerolog.NewConsoleWriter())

	if flag.Arg(0) == "cli" {
		if cliWebviewUrl == "" {
			println("webview url not set (--req)")
			os.Exit(1)
		}
		pr, err := server.Parse(cliWebviewUrl, serviceUrl)
		if err != nil {
			panic(err)
		}
		println(server.AllPlayerResume(pr).String())
		return
	}

	err := connection.InitConnection(connection.ConnectionSettings{DBUrl: databaseName, Token: tursoTokenValue})
	if err != nil {
		panic(err)
	}

	s := server.Server(serviceUrl, telegramTokenValue)
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
