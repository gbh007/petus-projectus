package migrator

// Диалект мигратора.
const (
	PostgreSQL = iota
	MariaDB
	ClickHouse
)

func getMigration(dialect int) (string, error) {
	switch dialect {
	case PostgreSQL:
		return techMigrationPostgreSQL, nil

	case MariaDB:
		return techMigrationMariaDB, nil

	default:
		return "", UnknownDialect
	}
}

// миграции для таблиц мигратора.
const (
	techMigrationPostgreSQL = `
CREATE TABLE IF NOT EXISTS migrations(
    id          INT4        PRIMARY KEY,
    filename    TEXT        NOT NULL,
    hash        TEXT        NOT NULL,
    applied     TIMESTAMPTZ NOT NULL
);
`
	techMigrationMariaDB = `
CREATE TABLE IF NOT EXISTS migrations(
    id          INT         PRIMARY KEY,
    filename    TEXT        NOT NULL,
    hash        TINYTEXT    NOT NULL,
    applied     TIMESTAMP  NOT NULL
);
`
)
