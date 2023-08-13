package client

import (
	"app/services/notification/internal/pb"
	"context"
	"errors"

	"google.golang.org/grpc"
)

type Client struct {
	client pb.NotificationClient
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
	c.client = pb.NewNotificationClient(conn)

	return c, nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return errors.New("no connection")
	}

	return c.conn.Close()
}

func (c *Client) New(ctx context.Context, userID int64, n *Notification) error {
	res, err := c.client.New(ctx, &pb.NewRequest{
		UserID: userID,
		Data: &pb.NotificationData{
			Kind:  n.Kind,
			Level: n.Level,
			Title: n.Title,
			Body:  n.Body,
		},
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
