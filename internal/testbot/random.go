package testbot

import (
	"crypto/sha256"
	"fmt"
	"time"
)

func randomSHA256String(prefix string) string {
	return fmt.Sprintf("%s%x", prefix, sha256.Sum256([]byte(time.Now().String())))
}
