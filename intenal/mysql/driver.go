package mysql

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/betterde/mysql-mcp-server/config"
	"github.com/betterde/mysql-mcp-server/intenal/journal"
)

var Conn *sql.Conn

func Init(ctx context.Context, conf *config.Config) {
	db, err := sql.Open("mysql", conf.DSN)
	if err != nil {
		journal.Logger.Panic(err.Error())
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 3)

	Conn, err = db.Conn(ctx)
	if err != nil {
		journal.Logger.Panic(err.Error())
	}
}
