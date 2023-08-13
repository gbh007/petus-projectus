package storage

import (
	"app/internal/migrator"
	"app/services/log/internal/storage/migration"
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mailru/go-clickhouse/v2"
)

var databaseErr = errors.New("log database")

type Database struct {
	db *sqlx.DB
}

func Init(ctx context.Context, username, password, dbHostWithPort, databaseName string) (*Database, error) {
	cs := fmt.Sprintf("http://%s:%s@%s/%s", username, password, dbHostWithPort, databaseName)

	db, err := sqlx.Open("chhttp", cs)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", databaseErr, err)
	}

	err = migrator.MigrateAll(ctx, migration.Migrations, db, true, migrator.ClickHouse)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", databaseErr, err)
	}

	return &Database{
		db: db,
	}, nil
}
