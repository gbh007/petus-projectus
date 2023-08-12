package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type UserLog struct {
	RequestID    string         `db:"request_id"`
	Addr         string         `db:"addr"`
	UserID       sql.NullInt64  `db:"user_id"`
	SessionToken sql.NullString `db:"session_token"`
	Action       string         `db:"action"`
	Chance       sql.NullInt64  `db:"chance"`
	Duration     sql.NullInt64  `db:"duration"`
	RequestTime  time.Time      `db:"request_time"`
}

func (db *Database) InsertUserLog(ctx context.Context, ul *UserLog) error {
	ul.RequestTime = ul.RequestTime.UTC()

	_, err := db.db.NamedExecContext(ctx, `INSERT INTO user_logs (
        request_id,
        addr,
        user_id,
        session_token,
        action,
        chance,
        duration,
        request_time
) VALUES (
        :request_id,
        :addr,
        :user_id,
        :session_token,
        :action,
        :chance,
        :duration,
        :request_time
);`, ul)
	if err != nil {
		return fmt.Errorf("%w: %w", databaseErr, err)
	}

	return nil
}
