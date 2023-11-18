package storage

import (
	"app/internal/migrator"
	"app/services/notification/internal/storage/migration"
	"context"
	"errors"
	"fmt"

	// import sql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
)

var errDatabase = errors.New("notification database")

type Database struct {
	db *sqlx.DB
}

func Init(ctx context.Context, username, password, dbHostWithPort, databaseName string) (*Database, error) {
	cs := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true&multiStatements=true",
		username, password, dbHostWithPort, databaseName,
	)

	db, err := otelsqlx.Open("mysql", cs)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errDatabase, err)
	}

	err = migrator.MigrateAll(ctx, migration.Migrations, db, true, migrator.MySQL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errDatabase, err)
	}

	return &Database{
		db: db,
	}, nil
}
