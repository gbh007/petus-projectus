package gatedto

import "time"

const (
	ActionLogin    = "login"
	ActionRegister = "register"
	ActionButton   = "button"
)

type KafkaData struct {
	Addr         string    `json:"addr,omitempty"`
	UserID       int64     `json:"user_id,omitempty"`
	SessionToken string    `json:"session_token,omitempty"`
	Action       string    `json:"action,omitempty"`
	Chance       int64     `json:"chance,omitempty"`
	Duration     int64     `json:"duration,omitempty"`
	RequestTime  time.Time `json:"request_time,omitempty"`
}
