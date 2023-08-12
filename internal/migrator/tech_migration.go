package migrator

// Диалект мигратора.
const (
	PostgreSQL = iota
	MySQL
	ClickHouse
)

func getMigration(dialect int) (string, error) {
	switch dialect {
	case PostgreSQL:
		return techMigrationPostgreSQL, nil

	case MySQL:
		return techMigrationMariaDB, nil

	case ClickHouse:
		return techMigrationClickHouse, nil

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
    hash        TEXT        NOT NULL,
    applied     TIMESTAMP   NOT NULL
);
`

	// Данное решение прототип.
	techMigrationClickHouse = `
CREATE TABLE IF NOT EXISTS migrations(
    id          Int64         NOT NULL,
    filename    String        NOT NULL,
    hash        String        NOT NULL,
    applied     DateTime64(9, 'UTC')   NOT NULL
)
ENGINE = MergeTree()
ORDER BY id;
`
)
