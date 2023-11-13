package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/logger"
)

var postgresConn *pgxpool.Pool
var postgresLog logger.Logger

func InitPostgres(ctx context.Context) {
	postgresLog = logger.Get("postgres")
	postgresLog.Info("Connect Postgres server...")
	conf := config.Get()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", conf.DB.User, conf.DB.Password, conf.DB.Host, conf.DB.Port, conf.DB.Name)
	connConf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		postgresLog.Fatal(err, "Postgres configration problem")
	}

	connConf.MaxConnIdleTime = conf.DB.ConnectionIdle
	connConf.MaxConnLifetime = conf.DB.ConnectionLifetime
	connConf.MaxConns = int32(conf.DB.MaxOpen)
	connConf.MinConns = int32(conf.DB.MaxIdle)

	db, err := pgxpool.NewWithConfig(ctx, connConf)
	if err != nil {
		postgresLog.Fatal(err, "Postgres connection problem")
	}

	defer db.Close()
	postgresConn = new(pgxpool.Pool)
	postgresConn = db

	postgresLog.Info("Postgres connected...")
}

func GetPostgres() *pgxpool.Pool {
	return postgresConn
}
