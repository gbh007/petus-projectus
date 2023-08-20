package desktop

import (
	"context"
	"time"
)

var _ Controller = new(controllerMock)

type controllerMock struct{}

func (c *controllerMock) Login(_ context.Context, _, _ string) error {
	return nil
}

func (c *controllerMock) Register(_ context.Context, _, _ string) error {
	return nil
}

func (c *controllerMock) ButtonClick(_ context.Context, _, _ int64) error {
	return nil
}

func (c *controllerMock) Notifications(_ context.Context) ([]Notification, error) {
	return []Notification{}, nil
}

func (c *controllerMock) Activity(_ context.Context) (int64, time.Time, error) {
	return 0, time.Time{}, nil
}
