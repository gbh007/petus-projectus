package client

import (
	"app/services/gate/internal/pb"
	"context"
	"errors"

	"google.golang.org/grpc"
)

type Client struct {
	gateClient         pb.GateClient
	notificationClient pb.NotificationClient
	conn               *grpc.ClientConn
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
	c.gateClient = pb.NewGateClient(conn)
	c.notificationClient = pb.NewNotificationClient(conn)

	return c, nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return errors.New("no connection")
	}

	return c.conn.Close()
}

func (c *Client) Login(ctx context.Context, login, pass string) (string, error) {
	res, err := c.gateClient.Login(ctx, &pb.LoginRequest{
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
	res, err := c.gateClient.Register(ctx, &pb.RegisterRequest{
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

func (c *Client) ButtonClick(ctx context.Context, token string, duration, chance int64) error {
	res, err := c.gateClient.Button(ctx, &pb.ButtonRequest{
		Duration: duration,
		Token:    token,
		Chance:   chance,
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

func (c *Client) Read(ctx context.Context, token string, all bool, id int64) error {
	res, err := c.notificationClient.Read(ctx, &pb.NotificationReadRequest{
		Token: token,
		Id:    id,
		All:   all,
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

func (c *Client) List(ctx context.Context, token string) ([]*Notification, error) {
	res, err := c.notificationClient.List(ctx, &pb.NotificationListRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	if res.GetError() != nil {
		err := errors.New(res.GetError().GetText())

		return nil, err
	}

	notifications := make([]*Notification, len(res.GetList()))

	for index, raw := range res.GetList() {
		notifications[index] = &Notification{
			ID:      raw.GetId(),
			Kind:    raw.GetKind(),
			Level:   raw.GetLevel(),
			Title:   raw.GetTitle(),
			Body:    raw.GetBody(),
			Created: raw.GetCreated().AsTime(),
		}
	}

	return notifications, nil
}
