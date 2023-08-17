package desktop

import (
	"context"
	"errors"
	"time"
)

var _ Controller = new(controllerMock)

type controllerMock struct{}

func (c *controllerMock) Login(ctx context.Context, login, pass string) error {
	if login == "err" {
		return errors.New(pass)
	}

	return nil
}

func (c *controllerMock) Register(ctx context.Context, login, pass string) error {
	if login == "err" {
		return errors.New(pass)
	}

	return nil
}

func (c *controllerMock) ButtonClick(ctx context.Context, duration, chance int64) error {
	if duration < 0 {
		return errors.New("duration less 0")
	}

	return nil
}

func (c *controllerMock) Notifications(ctx context.Context) ([]Notification, error) {
	return []Notification{}, nil
}

func (c *controllerMock) Activity(ctx context.Context) (int64, time.Time, error) {
	return 0, time.Time{}, nil
}
