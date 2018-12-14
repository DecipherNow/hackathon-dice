package github

import "time"

// Organization represents a GitHub organization account.
type Organization struct {
	Login                       string    `json:"login,omitempty"`
	ID                          int64     `json:"id,omitempty"`
	NodeID                      string    `json:"node_id,omitempty"`
	AvatarURL                   string    `json:"avatar_url,omitempty"`
	HTMLURL                     string    `json:"html_url,omitempty"`
	Name                        string    `json:"name,omitempty"`
	Company                     string    `json:"company,omitempty"`
	Blog                        string    `json:"blog,omitempty"`
	Location                    string    `json:"location,omitempty"`
	Email                       string    `json:"email,omitempty"`
	Description                 string    `json:"description,omitempty"`
	PublicRepos                 int       `json:"public_repos,omitempty"`
	PublicGists                 int       `json:"public_gists,omitempty"`
	Followers                   int       `json:"followers,omitempty"`
	Following                   int       `json:"following,omitempty"`
	CreatedAt                   time.Time `json:"created_at,omitempty"`
	UpdatedAt                   time.Time `json:"updated_at,omitempty"`
	TotalPrivateRepos           int       `json:"total_private_repos,omitempty"`
	OwnedPrivateRepos           int       `json:"owned_private_repos,omitempty"`
	PrivateGists                int       `json:"private_gists,omitempty"`
	DiskUsage                   int       `json:"disk_usage,omitempty"`
	Collaborators               int       `json:"collaborators,omitempty"`
	BillingEmail                string    `json:"billing_email,omitempty"`
	Type                        string    `json:"type,omitempty"`
	TwoFactorRequirementEnabled bool      `json:"two_factor_requirement_enabled,omitempty"`

	// DefaultRepoPermission can be one of: "read", "write", "admin", or "none". (Default: "read").
	// It is only used in OrganizationsService.Edit.
	DefaultRepoPermission string `json:"default_repository_permission,omitempty"`
	// DefaultRepoSettings can be one of: "read", "write", "admin", or "none". (Default: "read").
	// It is only used in OrganizationsService.Get.
	DefaultRepoSettings string `json:"default_repository_settings,omitempty"`

	// MembersCanCreateRepos default value is true and is only used in Organizations.Edit.
	MembersCanCreateRepos bool `json:"members_can_create_repositories,omitempty"`

	// API URLs
	URL              string `json:"url,omitempty"`
	EventsURL        string `json:"events_url,omitempty"`
	HooksURL         string `json:"hooks_url,omitempty"`
	IssuesURL        string `json:"issues_url,omitempty"`
	MembersURL       string `json:"members_url,omitempty"`
	PublicMembersURL string `json:"public_members_url,omitempty"`
	ReposURL         string `json:"repos_url,omitempty"`
}
