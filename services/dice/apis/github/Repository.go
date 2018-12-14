package github

import "time"

// Repository represents a GitHub repository.
type Repository struct {
	ID               int64     `json:"id,omitempty"`
	NodeID           string    `json:"node_id,omitempty"`
	Name             string    `json:"name,omitempty"`
	FullName         string    `json:"full_name,omitempty"`
	Description      string    `json:"description,omitempty"`
	Language         string    `json:"language,omitempty"`
	DefaultBranch    string    `json:"default_branch,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	PushedAt         time.Time `json:"pushed_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
	Fork             bool      `json:"fork,omitempty"`
	Private          bool      `json:"private,omitempty"`
	Archived         bool      `json:"archived,omitempty"`
	ForksCount       int       `json:"forks_count,omitempty"`
	NetworkCount     int       `json:"network_count,omitempty"`
	OpenIssuesCount  int       `json:"open_issues_count,omitempty"`
	StargazersCount  int       `json:"stargazers_count,omitempty"`
	SubscribersCount int       `json:"subscribers_count,omitempty"`
	WatchersCount    int       `json:"watchers_count,omitempty"`
	Size             int       `json:"size,omitempty"`
}
