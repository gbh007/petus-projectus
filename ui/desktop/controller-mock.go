package desktop

import (
	"context"
	"errors"
)

type ControllerMock struct{}

func (c *ControllerMock) Login(ctx context.Context, login, pass string) error {
	if login == "err" {
		return errors.New(pass)
	}

	return nil
}

func (c *ControllerMock) Register(ctx context.Context, login, pass string) error {
	if login == "err" {
		return errors.New(pass)
	}

	return nil
}

func (c *ControllerMock) ButtonClick(ctx context.Context, duration, chance int64) error {
	if duration < 0 {
		return errors.New("duration less 0")
	}

	return nil
}

func (c *ControllerMock) Notifications(ctx context.Context) ([]Notification, error) {
	return []Notification{}, nil
}
