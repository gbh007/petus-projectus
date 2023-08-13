package desktop

import (
	gateClient "app/services/gate/client"
	"context"
	"errors"
)

type controller_gRPC struct {
	client *gateClient.Client
	token  string
}

func newController_gRPC(addr string) (*controller_gRPC, error) {
	c := new(controller_gRPC)

	gc, err := gateClient.New(addr)
	if err != nil {
		return nil, err
	}

	c.client = gc

	return c, nil
}

func (c *controller_gRPC) Close() error {
	if c.client == nil {
		return errors.New("no client")
	}

	return c.client.Close()
}

func (c *controller_gRPC) Login(ctx context.Context, login, pass string) error {
	token, err := c.client.Login(ctx, login, pass)
	if err != nil {
		return err
	}

	c.token = token

	return nil
}

func (c *controller_gRPC) Register(ctx context.Context, login, pass string) error {
	err := c.client.Register(ctx, login, pass)
	if err != nil {
		return err
	}

	return nil
}

func (c *controller_gRPC) ButtonClick(ctx context.Context, duration, chance int64) error {
	err := c.client.ButtonClick(ctx, c.token, duration, chance)
	if err != nil {
		return err
	}

	return nil
}

func (c *controller_gRPC) Notifications(ctx context.Context) ([]Notification, error) {
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
