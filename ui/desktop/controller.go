package desktop

import (
	"context"
	"time"
)

type Controller interface {
	Login(ctx context.Context, login, pass string) error
	Register(ctx context.Context, login, pass string) error
	ButtonClick(ctx context.Context, duration, chance int64) error
	Notifications(ctx context.Context) ([]Notification, error)
}

type Notification struct {
	ID      int64
	IsOK    bool
	Title   string
	Body    string
	Created time.Time
}
