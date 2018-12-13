package github

import (
	"time"
)

type User struct {
	Login     string    `json:"login"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Created   time.Time `json:"created_at"`
	Updated   time.Time `json:"updated_at"`
}
