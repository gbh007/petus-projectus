package migrator

import (
	"context"
	"fmt"
	"io/fs"
	"log"

	"github.com/jmoiron/sqlx"
)

// MigrateAll - производит накат всех доступных миграций
func MigrateAll(ctx context.Context, migrationsDir fs.FS, db *sqlx.DB, checkHash bool, dialect int) error {
	list, err := getFileList(ctx, migrationsDir)
	if err != nil {
		return fmt.Errorf("%w: %w", MigratorError, err)
	}

	log.Println("Доступные миграции")
	for _, item := range list {
		log.Printf("%4d > %s\n", item.Number, item.Name)
	}

	var success bool

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", MigratorError, err)
	}

	// Функция для финализации транзакции
	defer func() {
		if success {
			logIfErr(tx.Commit())
		} else {
			logIfErr(tx.Rollback())
		}
	}()

	techMigration, err := getMigration(dialect)
	if err != nil {
		return fmt.Errorf("%w: %w", MigratorError, err)
	}

	_, err = tx.ExecContext(ctx, techMigration)
	if err != nil {
		return fmt.Errorf("%w: %w", MigratorError, err)
	}

	// Получаем список примененных миграций из БД
	appliedMigrationsList, err := getAppliedMigration(ctx, tx)
	if err != nil {
		return fmt.Errorf("%w: %w", MigratorError, err)
	}

	appliedMigrationsMap := make(map[int]mteMigration)
	for _, migration := range appliedMigrationsList {
		appliedMigrationsMap[migration.ID] = migration
	}

	for _, migrationFile := range list {
		mteInfo, migrationApplied := appliedMigrationsMap[migrationFile.Number]

		var (
			hash, body string
		)

		// Не применена миграция или нужно сверить ее хеш
		if !migrationApplied || checkHash {
			body, hash, err = migrationFromFile(ctx, migrationFile, migrationsDir)
			if err != nil {
				return fmt.Errorf("%w: %w", MigratorError, err)
			}
		}

		hashEqual := mteInfo.Hash == hash

		switch {
		// Миграция не применена
		case !migrationApplied:
			err = applyMigration(ctx, tx, migrationFile.Number, body, migrationFile.Name, hash)
			if err != nil {
				log.Printf("%s - ERR\n", migrationFile.Name)

				return fmt.Errorf("%w: %w", MigratorError, err)
			}

			log.Printf("%s - OK\n", migrationFile.Name)

		// Миграция применена и хеш не эквивалентен
		case migrationApplied && checkHash && !hashEqual:
			log.Println("old - hash >>>", mteInfo.Hash)
			log.Println("new - hash >>>", hash)
			log.Printf("%s - Not Equal HASH\n", migrationFile.Name)

		// Миграция уже применена
		case migrationApplied:
			log.Printf("%s - EXIST %v\n", migrationFile.Name, mteInfo.Applied)

		}
	}

	success = true

	return nil
}
