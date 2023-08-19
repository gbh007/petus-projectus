package client

import (
	"app/services/log/internal/pb"
	"context"
	"errors"

	"google.golang.org/grpc"
)

type Client struct {
	client pb.LogClient
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
	c.client = pb.NewLogClient(conn)

	return c, nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return errors.New("no connection")
	}

	return c.conn.Close()
}

func (c *Client) Activity(ctx context.Context, userID int64) (*LogData, error) {
	res, err := c.client.Activity(ctx, &pb.ActivityRequest{
		UserID: userID,
	})

	if err != nil {
		return nil, err
	}

	return &LogData{
		RequestCount: res.GetData().GetRequestCount(),
		LastRequest:  res.GetData().GetLastRequest().AsTime(),
	}, nil
}
