package github

type Member struct {
	Login      string `json:"login"`
	ID         int    `json:"id"`
	NodeID     string `json:"node_id"`
	AvatarURL  string `json:"avatar_url,omitempty"`
	GravatarID string `json:"gavatar_id,omitempty"`
	URL        string `json:"url,omitempty"`
	Type       string `json:"type"`
	SiteAdmin  bool   `json:"site_admin"`
}

type MembersResponse struct {
	Members []Member
}
