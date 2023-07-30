package gate

import (
	"app/internal/auth/authpb"
	"context"
	"errors"

	"google.golang.org/grpc"
)

type authClient struct {
	client authpb.AuthClient
	conn   *grpc.ClientConn
}

func newAuthClient(addr string) (*authClient, error) {
	c := new(authClient)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c.conn = conn
	c.client = authpb.NewAuthClient(conn)

	return c, nil
}

func (c *authClient) Close() error {
	if c.conn == nil {
		return errors.New("no connection")
	}

	return c.conn.Close()
}

func (c *authClient) Login(ctx context.Context, login, pass string) (string, error) {
	res, err := c.client.Login(ctx, &authpb.LoginRequest{
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

func (c *authClient) Register(ctx context.Context, login, pass string) error {
	res, err := c.client.Register(ctx, &authpb.RegisterRequest{
		Login:    login,
		Password: pass,
	})

	if err != nil {
		return err
	}

	if res.GetError().GetHas() {
		// code := res.GetError().GetCode()
		err := errors.New(res.GetError().GetText())

		return err
	}

	return nil
}
