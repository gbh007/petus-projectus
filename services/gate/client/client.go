package client

import (
	"app/services/gate/internal/gatepb"
	"context"
	"errors"

	"google.golang.org/grpc"
)

type Client struct {
	client gatepb.GateClient
	conn   *grpc.ClientConn
}

type UserInfo struct {
	ID int64
}

func New(addr string) (*Client, error) {
	c := new(Client)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c.conn = conn
	c.client = gatepb.NewGateClient(conn)

	return c, nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return errors.New("no connection")
	}

	return c.conn.Close()
}

func (c *Client) Login(ctx context.Context, login, pass string) (string, error) {
	res, err := c.client.Login(ctx, &gatepb.LoginRequest{
		Login:    login,
		Password: pass,
	})

	if err != nil {
		return "", err
	}

	if res.GetError().GetHas() {
		err := errors.New(res.GetError().GetText())

		return "", err
	}

	return res.GetToken(), nil
}

func (c *Client) Register(ctx context.Context, login, pass string) error {
	res, err := c.client.Register(ctx, &gatepb.RegisterRequest{
		Login:    login,
		Password: pass,
	})

	if err != nil {
		return err
	}

	if res.GetError().GetHas() {
		err := errors.New(res.GetError().GetText())

		return err
	}

	return nil
}

func (c *Client) ButtonClick(ctx context.Context, token string, duration, chance int64) error {
	res, err := c.client.Button(ctx, &gatepb.ButtonRequest{
		Duration: duration,
		Token:    token,
		Chance:   chance,
	})

	if err != nil {
		return err
	}

	if res.GetError().GetHas() {
		err := errors.New(res.GetError().GetText())

		return err
	}

	return nil
}
