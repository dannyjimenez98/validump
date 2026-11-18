package postgres

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dannyjimenez98/validump/internal/config"
	"github.com/dannyjimenez98/validump/internal/terminal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spf13/viper"
)

func Connect(ctx context.Context, cfg *viper.Viper, dbalias string) (*pgx.Conn, error) {
	dbCfg, err := config.LoadDB(cfg, dbalias)
	if err != nil {
		return nil, fmt.Errorf("get config for database '%s': %w", dbalias, err)
	}

	pgxCfg, err := NewPgxConfig(dbalias, dbCfg)
	if err != nil {
		return nil, fmt.Errorf("create PGX configuration: %w", err)
	}

	connection, err := connectWithTimeout(ctx, pgxCfg)
	if err == nil {
		return connection, nil
	}

	var pgErr *pgconn.PgError

	// check if err is not authentication related
	if pgxCfg.Password != "" || !errors.As(err, &pgErr) || pgErr.Code != "28P01" {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	password, err := terminal.GetDBPassword()
	if err != nil {
		return nil, err
	}
	defer clear(password)

	updatedCfg := pgxCfg.Copy()
	updatedCfg.Password = string(password)

	connection, err = connectWithTimeout(ctx, updatedCfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return connection, nil
}

func NewPgxConfig(alias string, db config.DBConfig) (*pgx.ConnConfig, error) {
	uri := url.URL{
		Scheme: "postgres",
		User:   url.User(db.User),
		Host:   net.JoinHostPort(db.Host, strconv.Itoa(int(db.Port))),
		Path:   "/" + db.DBName,
	}

	pgxCfg, err := pgx.ParseConfig(uri.String())
	if err != nil {
		return nil, fmt.Errorf("parse pgx config: %w", err)
	}

	envPassword := fmt.Sprintf("VALIDUMP_DATABASES__%s__PASSWORD", strings.ToUpper(alias))
	if pw := os.Getenv(envPassword); pw != "" {
		pgxCfg.Password = pw
	} else if db.Password != "" {
		pgxCfg.Password = db.Password
	}

	return pgxCfg, nil
}

func connectWithTimeout(ctx context.Context, cfg *pgx.ConnConfig) (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return pgx.ConnectConfig(ctx, cfg)
}
