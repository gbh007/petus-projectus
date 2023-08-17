package client

import (
	"app/services/auth/internal/pb"
	"context"
	"errors"

	"google.golang.org/grpc"
)

type Client struct {
	client pb.AuthClient
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
	c.client = pb.NewAuthClient(conn)

	return c, nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return errors.New("no connection")
	}

	return c.conn.Close()
}

func (c *Client) Login(ctx context.Context, login, pass string) (string, error) {
	res, err := c.client.Login(ctx, &pb.LoginRequest{
		Login:    login,
		Password: pass,
	})

	if err != nil {
		return "", err
	}

	if res.GetError() != nil {
		err := errors.New(res.GetError().GetText())

		return "", err
	}

	return res.GetToken(), nil
}

func (c *Client) Register(ctx context.Context, login, pass string) error {
	res, err := c.client.Register(ctx, &pb.RegisterRequest{
		Login:    login,
		Password: pass,
	})

	if err != nil {
		return err
	}

	if res.GetError() != nil {
		err := errors.New(res.GetError().GetText())

		return err
	}

	return nil
}

func (c *Client) Info(ctx context.Context, token string) (*UserInfo, error) {
	res, err := c.client.Info(ctx, &pb.InfoRequest{
		Token: token,
	})

	if err != nil {
		return nil, err
	}

	if res.GetError() != nil {
		err := errors.New(res.GetError().GetText())

		return nil, err
	}

	return &UserInfo{
		ID: res.GetUserID(),
	}, nil
}
