package client

const (
	ButtonKind   = "button"
	SuccessLevel = "success"
	ErrorLevel   = "error"
)

type Notification struct {
	Kind  string
	Level string
	Title string
	Body  string
}
