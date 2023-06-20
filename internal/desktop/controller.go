package desktop

import "errors"

type Controller interface {
	Login(login, pass string) error
	Register(login, pass string) error
	ButtonClick(duration, chance int64) error
}

type ControllerMock struct{}

func (c *ControllerMock) Login(login, pass string) error {
	if login == "err" {
		return errors.New(pass)
	}

	return nil
}

func (c *ControllerMock) Register(login, pass string) error {
	if login == "err" {
		return errors.New(pass)
	}

	return nil
}

func (c *ControllerMock) ButtonClick(duration, chance int64) error {
	if duration < 0 {
		return errors.New("duration less 0")
	}

	return nil
}
