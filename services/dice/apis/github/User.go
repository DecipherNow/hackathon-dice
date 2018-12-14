package github

type User struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Created   string `json:"created_at"`
	Updated   string `json:"updated_at"`
}
