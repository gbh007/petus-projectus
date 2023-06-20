package desktop

import (
	"app/internal/gate/gatepb"
	"context"
	"errors"

	"google.golang.org/grpc"
)

type controller_gRPC struct {
	client gatepb.GateClient
	conn   *grpc.ClientConn
	token  string
}

func newController_gRPC(addr string) (*controller_gRPC, error) {
	c := new(controller_gRPC)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c.conn = conn
	c.client = gatepb.NewGateClient(conn)

	return c, nil
}

func (c *controller_gRPC) Close() error {
	if c.conn == nil {
		return errors.New("no connection")
	}

	return c.conn.Close()
}

func (c *controller_gRPC) Login(ctx context.Context, login, pass string) error {
	res, err := c.client.Login(ctx, &gatepb.LoginRequest{
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

	c.token = res.GetToken()

	return nil
}

func (c *controller_gRPC) Register(ctx context.Context, login, pass string) error {
	res, err := c.client.Register(ctx, &gatepb.RegisterRequest{
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

func (c *controller_gRPC) ButtonClick(ctx context.Context, duration, chance int64) error {
	res, err := c.client.Button(ctx, &gatepb.ButtonRequest{
		Duration: duration,
		Token:    c.token,
		Chance:   chance,
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
