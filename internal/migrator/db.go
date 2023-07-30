package migrator

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

// mteMigration - модель для таблицы с миграциями
type mteMigration struct {
	ID       int       `db:"id"`
	Filename string    `db:"filename"`
	Hash     string    `db:"hash"`
	Applied  time.Time `db:"applied"`
}

// getAppliedMigration - возвращает примененные миграции, упорядоченные по номеру
func getAppliedMigration(ctx context.Context, tx *sqlx.Tx) ([]mteMigration, error) {
	migrations := make([]mteMigration, 0)

	err := tx.SelectContext(ctx, &migrations, `SELECT * FROM migrations ORDER BY id;`)
	if err != nil {
		return nil, err
	}

	return migrations, nil
}

// applyMigration - применяет миграцию
func applyMigration(ctx context.Context, tx *sqlx.Tx, id int, body, filename, hash string) error {
	_, err := tx.ExecContext(ctx, body)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO migrations (id, filename, hash, applied) VALUES ($1, $2, $3, $4);`,
		id, filename, hash, time.Now().UTC(),
	)
	if err != nil {
		return err
	}

	return nil
}
