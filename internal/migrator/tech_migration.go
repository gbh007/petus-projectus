package migrator

// techMigrationPostgreSQL - миграция для таблиц мигратора
const techMigrationPostgreSQL = `
CREATE TABLE IF NOT EXISTS migrations(
    id          INT4        PRIMARY KEY,
    filename    TEXT        NOT NULL,
    hash        TEXT        NOT NULL,
    applied     TIMESTAMPTZ NOT NULL
);
`
