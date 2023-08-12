package storage

import (
	"app/internal/migrator"
	"app/services/auth/internal/storage/migration"
	"context"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var databaseErr = errors.New("auth database")

type Database struct {
	db *sqlx.DB
}

func Init(ctx context.Context, username, password, dbHostWithPort, databaseName string) (*Database, error) {
	cs := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&multiStatements=true", username, password, dbHostWithPort, databaseName)

	db, err := sqlx.Open("mysql", cs)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", databaseErr, err)
	}

	err = migrator.MigrateAll(ctx, migration.Migrations, db, true, migrator.MySQL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", databaseErr, err)
	}

	return &Database{
		db: db,
	}, nil
}
