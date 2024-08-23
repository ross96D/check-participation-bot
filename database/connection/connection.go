package connection

import (
	"database/sql"
	"net/url"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type ConnectionSettings struct {
	DBUrl string
	Token string
}

var db *sql.DB

func DB() *sql.DB {
	return db
}

func InitConnection(settings ConnectionSettings) (err error) {
	uri, err := url.Parse(settings.DBUrl)
	if err != nil {
		uri, err = url.Parse("libsql://" + settings.DBUrl)
		if err != nil {
			return err
		}
	}

	if settings.Token != "" {
		q := uri.Query()
		q.Add("authToken", settings.Token)
		uri.RawQuery = q.Encode()
	}
	db, err = sql.Open("libsql", uri.String())
	return err
}
