package desktop

import (
	gateClient "app/services/gate/client"
	"context"
	"errors"
	"time"
)

var _ Controller = new(controllerGate)

var errNoClient = errors.New("no client")

type controllerGate struct {
	client *gateClient.Client
	token  string
}

func newControllerGate(addr string) (*controllerGate, error) {
	c := new(controllerGate)

	gc, err := gateClient.New(addr)
	if err != nil {
		return nil, err
	}

	c.client = gc

	return c, nil
}

func (c *controllerGate) Close() error {
	if c.client == nil {
		return errNoClient
	}

	return c.client.Close()
}

func (c *controllerGate) Login(ctx context.Context, login, pass string) error {
	token, err := c.client.Login(ctx, login, pass)
	if err != nil {
		return err
	}

	c.token = token

	return nil
}

func (c *controllerGate) Register(ctx context.Context, login, pass string) error {
	err := c.client.Register(ctx, login, pass)
	if err != nil {
		return err
	}

	return nil
}

func (c *controllerGate) ButtonClick(ctx context.Context, duration, chance int64) error {
	err := c.client.ButtonClick(ctx, c.token, duration, chance)
	if err != nil {
		return err
	}

	return nil
}

func (c *controllerGate) Notifications(ctx context.Context) ([]Notification, error) {
	data, err := c.client.List(ctx, c.token)
	if err != nil {
		return nil, err
	}

	out := make([]Notification, len(data))
	for i, raw := range data {
		out[i] = Notification{
			ID:      raw.ID,
			IsOK:    raw.Level == gateClient.SuccessLevel,
			Title:   raw.Title,
			Body:    raw.Body,
			Created: raw.Created,
		}
	}

	return out, nil
}

func (c *controllerGate) Activity(ctx context.Context) (int64, time.Time, error) {
	return c.client.Activity(ctx, c.token)
}
