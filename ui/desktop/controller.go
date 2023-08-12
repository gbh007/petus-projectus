package desktop

import "context"

type Controller interface {
	Login(ctx context.Context, login, pass string) error
	Register(ctx context.Context, login, pass string) error
	ButtonClick(ctx context.Context, duration, chance int64) error
}
