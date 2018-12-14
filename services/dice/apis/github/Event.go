package github

type Actor struct {
	Login string `json:"login"`
}
type Author struct {
	Name string `json:"name"`
}
type Comment struct {
	DiffHunk  string `json:"diff_hunk"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
}
type Commit struct {
	Author  Author `json:"author"`
	Message string `json:"message"`
}
type Event struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"created_at"`
	Actor     Actor   `json:"actor"`
	Payload   Payload `json:"payload"`
	Repo      Repo    `json:"repo"`
	Type      string  `json:"type"`
}
type Issue struct {
	Number      int         `json:"number"`
	Title       string      `json:"title"`
	Comments    int         `json:"comments"`
	Body        string      `json:"Body"`
	PullRequest PullRequest `json:"pull_request"`
}
type Payload struct {
	Action      string      `json:"action"`
	Comment     Comment     `json:"comment"`
	Commits     []Commit    `json:"commits"`
	Issue       Issue       `json:"issue"`
	Number      int         `json:"number"`
	PullRequest PullRequest `json:"pull_request"`
	Ref         string      `json:"ref"`
	RefType     string      `json:"ref_type"`
	Size        int         `json:"size"`
}
type PullRequest struct {
	URI       string `json:"uri"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	MergedAt  string `json:"merged_at"`
	Merged    bool   `json:"merged"`
}
type Repo struct {
	Name string `json:"name"`
}
